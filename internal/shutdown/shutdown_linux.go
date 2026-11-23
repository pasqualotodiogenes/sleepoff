//go:build linux

package shutdown

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"os/exec"
	"time"
)

//go:embed countdown.wav
var countdownWavData []byte

// Execute realiza o desligamento real no Linux via systemctl poweroff,
// com fallback para shutdown -h now e poweroff se sem systemd.
func Execute(dryRun bool) error {
	if dryRun {
		return nil
	}

	// 1. Tenta systemctl poweroff
	if err := exec.Command("systemctl", "poweroff").Run(); err == nil {
		return nil
	}

	// 2. Fallback: shutdown -h now
	if err := exec.Command("shutdown", "-h", "now").Run(); err == nil {
		return nil
	}

	// 3. Fallback: poweroff
	return exec.Command("poweroff").Run()
}

// Shutdown é o alias idiomático para Execute.
func Shutdown(dryRun bool) error {
	return Execute(dryRun)
}

// Beep emite som no Linux via paplay, pw-play ou aplay se disponíveis,
// com fallback para o terminal bell (\a).
func Beep(frequency, duration int) {
	go playLinuxBeep(frequency, duration)
}

func playLinuxBeep(frequency, durationMs int) {
	wavData := countdownWavData
	if len(wavData) == 0 {
		if frequency <= 0 {
			frequency = 800
		}
		if durationMs <= 0 {
			durationMs = 100
		}
		wavData = generateSineWav(frequency, time.Duration(durationMs)*time.Millisecond)
	}

	// 1. Tenta pw-play (PipeWire)
	if path, err := exec.LookPath("pw-play"); err == nil {
		cmd := exec.Command(path, "-")
		cmd.Stdin = bytes.NewReader(wavData)
		if cmd.Run() == nil {
			return
		}
	}

	// 2. Tenta paplay (PulseAudio)
	if path, err := exec.LookPath("paplay"); err == nil {
		tmpFile, err := os.CreateTemp("", "sleepoff-beep-*.wav")
		if err == nil {
			_, _ = tmpFile.Write(wavData)
			_ = tmpFile.Close()
			defer os.Remove(tmpFile.Name())

			cmd := exec.Command(path, tmpFile.Name())
			if cmd.Run() == nil {
				return
			}
		}
	}

	// 3. Tenta aplay (ALSA)
	if path, err := exec.LookPath("aplay"); err == nil {
		cmd := exec.Command(path, "-q")
		cmd.Stdin = bytes.NewReader(wavData)
		if cmd.Run() == nil {
			return
		}
	}

	// Fallback universal: terminal bell
	fmt.Fprint(os.Stderr, "\a")
}

// generateSineWav cria um buffer WAV mono 16-bit PCM puro sem dependências externas.
func generateSineWav(frequency int, duration time.Duration) []byte {
	const sampleRate = 8000
	numSamples := int(float64(sampleRate) * duration.Seconds())
	if numSamples <= 0 {
		numSamples = sampleRate / 10
	}

	dataSize := numSamples * 2
	buf := new(bytes.Buffer)

	// RIFF header
	buf.WriteString("RIFF")
	_ = binary.Write(buf, binary.LittleEndian, uint32(36+dataSize))
	buf.WriteString("WAVE")

	// fmt subchunk
	buf.WriteString("fmt ")
	_ = binary.Write(buf, binary.LittleEndian, uint32(16))           // Subchunk1Size (16 for PCM)
	_ = binary.Write(buf, binary.LittleEndian, uint16(1))            // AudioFormat (1 = PCM)
	_ = binary.Write(buf, binary.LittleEndian, uint16(1))            // NumChannels (1 = Mono)
	_ = binary.Write(buf, binary.LittleEndian, uint32(sampleRate))   // SampleRate
	_ = binary.Write(buf, binary.LittleEndian, uint32(sampleRate*2)) // ByteRate
	_ = binary.Write(buf, binary.LittleEndian, uint16(2))            // BlockAlign
	_ = binary.Write(buf, binary.LittleEndian, uint16(16))           // BitsPerSample

	// data subchunk
	buf.WriteString("data")
	_ = binary.Write(buf, binary.LittleEndian, uint32(dataSize))

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(sampleRate)
		sample := int16(math.Sin(2*math.Pi*float64(frequency)*t) * 28000)
		_ = binary.Write(buf, binary.LittleEndian, sample)
	}

	return buf.Bytes()
}
