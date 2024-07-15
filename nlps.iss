[Setup]
AppName=NLPS
AppVersion=1.0
DefaultDirName={commonpf}\NLPS
DefaultGroupName=NLPS
Compression=lzma2
SolidCompression=yes
OutputDir=bin
UninstallDisplayIcon={app}\nlps.exe
UninstallDisplayName=NLPS
ChangesEnvironment=yes

[Files]
Source: "nlps.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "README.md"; DestDir: "{app}"; Flags: ignoreversion isreadme

[Icons]
Name: "{group}\NLPS"; Filename: "{app}\nlps.exe"
Name: "{group}\Uninstall NLPS"; Filename: "{uninstallexe}"

[Registry]
Root: HKLM; Subkey: "SYSTEM\CurrentControlSet\Control\Session Manager\Environment"; ValueType: expandsz; ValueName: "Path"; ValueData: "{olddata};{app}"; Check: NeedsAddPath(ExpandConstant('{app}'));

[Code]
function NeedsAddPath(Param: string): boolean;
var
  OrigPath: string;
begin
  if not RegQueryStringValue(HKEY_LOCAL_MACHINE,
    'SYSTEM\CurrentControlSet\Control\Session Manager\Environment',
    'Path', OrigPath)
  then begin
    Result := True;
    exit;
  end;
  // Avoid duplicate entries in the path
  Result := Pos(';' + Uppercase(Param) + ';', ';' + Uppercase(OrigPath) + ';') = 0;
end;

function RemoveFromPath(Path: string; Dir: string): string;
var
  P: Integer;
begin
  Result := Path;
  P := Pos(';' + Uppercase(Dir) + ';', ';' + Uppercase(Result) + ';');
  while P > 0 do
  begin
    Delete(Result, P, Length(Dir) + 1);
    P := Pos(';' + Uppercase(Dir) + ';', ';' + Uppercase(Result) + ';');
  end;
  if Copy(Result, Length(Result), 1) = ';' then
    Delete(Result, Length(Result), 1);
  if Copy(Result, 1, 1) = ';' then
    Delete(Result, 1, 1);
end;

function GetUserConfigDir(): string;
begin
  Result := ExpandConstant('{userappdata}\nlps');
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
var
  Path: string;
  AppDir: string;
  NewPath: string;
  ConfigDir: string;
  ConfigFile: string;
begin
  if CurUninstallStep = usUninstall then
  begin
    if RegQueryStringValue(HKEY_LOCAL_MACHINE,
      'SYSTEM\CurrentControlSet\Control\Session Manager\Environment',
      'Path', Path) then
    begin
      AppDir := ExpandConstant('{app}');
      NewPath := RemoveFromPath(Path, AppDir);
      RegWriteExpandStringValue(HKEY_LOCAL_MACHINE,
        'SYSTEM\CurrentControlSet\Control\Session Manager\Environment',
        'Path', NewPath);
    end;

    // Remove configuration file
    ConfigDir := GetUserConfigDir();
    ConfigFile := ConfigDir + '\config.yaml';
    if FileExists(ConfigFile) then
    begin
      if not DeleteFile(ConfigFile) then
      begin
        MsgBox('Failed to delete config file: ' + ConfigFile + #13#10 +
               'You may need to delete it manually.', mbInformation, MB_OK);
      end;
    end;
    
    // Attempt to remove the config directory if it's empty
    if not RemoveDir(ConfigDir) then
    begin
      MsgBox('Could not remove config directory: ' + ConfigDir + #13#10 +
             'It may not be empty or you may need to remove it manually.', mbInformation, MB_OK);
    end;
  end;
end;