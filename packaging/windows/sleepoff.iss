#ifndef MyAppVersion
  #define MyAppVersion "0.0.0"
#endif

#ifndef MyAppExe
  #define MyAppExe "sleepoff.exe"
#endif

#ifndef MyOutputDir
  #define MyOutputDir "dist"
#endif

[Setup]
AppId={{7E5BCDF0-7B9A-4F30-8D0B-4E7E1E43C837}
AppName=sleepoff
AppVersion={#MyAppVersion}
AppPublisher=Diogenes Pasqualoto
AppPublisherURL=https://github.com/pasqualotodiogenes/sleepoff
AppSupportURL=https://github.com/pasqualotodiogenes/sleepoff/issues
AppUpdatesURL=https://github.com/pasqualotodiogenes/sleepoff/releases/latest
DefaultDirName={localappdata}\Programs\sleepoff
DefaultGroupName=sleepoff
DisableProgramGroupPage=yes
LicenseFile={#SourcePath}\..\..\LICENSE
OutputDir={#MyOutputDir}
OutputBaseFilename=sleepoff-setup
SetupIconFile={#SourcePath}\icon.ico
UninstallDisplayIcon={app}\sleepoff.exe
Compression=lzma
SolidCompression=yes
WizardStyle=modern
PrivilegesRequired=lowest
ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible
ChangesEnvironment=yes

[Tasks]
Name: addtopath; Description: "Add sleepoff to your PATH"; GroupDescription: "Command line integration:"; Flags: checkedonce
Name: startmenuicon; Description: "Create a Start Menu shortcut"; GroupDescription: "Shortcuts:"; Flags: unchecked

[Files]
Source: "{#MyAppExe}"; DestDir: "{app}"; DestName: "sleepoff.exe"; Flags: ignoreversion
Source: "{#SourcePath}\..\..\README.md"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#SourcePath}\..\..\LICENSE"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{userprograms}\sleepoff"; Filename: "{app}\sleepoff.exe"; WorkingDir: "{app}"; Tasks: startmenuicon

[Run]
Filename: "{app}\sleepoff.exe"; Description: "Launch sleepoff now"; Flags: nowait postinstall skipifsilent unchecked

[Code]
const
  EnvKey = 'Environment';
  EnvValue = 'Path';

function NormalizePath(const Value: string): string;
begin
  Result := RemoveBackslashUnlessRoot(Trim(Value));
end;

function PathContainsDir(const PathValue, Dir: string): Boolean;
var
  Parts: TArrayOfString;
  I: Integer;
begin
  Result := False;
  Parts := SplitString(PathValue, ';');
  for I := 0 to GetArrayLength(Parts) - 1 do
  begin
    if CompareText(NormalizePath(Parts[I]), NormalizePath(Dir)) = 0 then
    begin
      Result := True;
      Exit;
    end;
  end;
end;

function RemovePathEntry(const PathValue, Dir: string): string;
var
  Parts: TArrayOfString;
  I: Integer;
  Item: string;
begin
  Result := '';
  Parts := SplitString(PathValue, ';');
  for I := 0 to GetArrayLength(Parts) - 1 do
  begin
    Item := Trim(Parts[I]);
    if (Item <> '') and (CompareText(NormalizePath(Item), NormalizePath(Dir)) <> 0) then
    begin
      if Result <> '' then
        Result := Result + ';';
      Result := Result + Item;
    end;
  end;
end;

procedure UpdateUserPath(const AddDir: Boolean);
var
  PathValue: string;
  NewPath: string;
  AppDir: string;
begin
  AppDir := ExpandConstant('{app}');
  if not RegQueryStringValue(HKCU, EnvKey, EnvValue, PathValue) then
    PathValue := '';

  if AddDir then
  begin
    if not PathContainsDir(PathValue, AppDir) then
    begin
      if (PathValue <> '') and (Copy(PathValue, Length(PathValue), 1) <> ';') then
        PathValue := PathValue + ';';
      RegWriteExpandStringValue(HKCU, EnvKey, EnvValue, PathValue + AppDir);
    end;
  end
  else
  begin
    NewPath := RemovePathEntry(PathValue, AppDir);
    if NewPath = '' then
      RegDeleteValue(HKCU, EnvKey, EnvValue)
    else
      RegWriteExpandStringValue(HKCU, EnvKey, EnvValue, NewPath);
  end;
end;

procedure CurStepChanged(CurStep: TSetupStep);
begin
  if (CurStep = ssPostInstall) and WizardIsTaskSelected('addtopath') then
    UpdateUserPath(True);
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  if CurUninstallStep = usUninstall then
    UpdateUserPath(False);
end;
