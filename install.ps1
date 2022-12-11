
if (Get-Command choco -errorAction SilentlyContinue)
{
    Write-Host "choco installed, step skipped"
}
else {
    Write-Host "download choco"
   
    iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
}

if (Get-Command go -errorAction SilentlyContinue)
{
    Write-Host "golang installed, step skipped"
}
else {
    Write-Host "golang choco"
   
    choco install golang -y
}


if (Get-Command go -errorAction SilentlyContinue)
{
    Write-Host "task installed, step skipped"
}
else {
    Write-Host "task choco"
   
    choco install go-task -y
}


task  build-debug

#userprrofile
$bashyHome=$([Environment]::GetFolderPath("UserProfile"))+"\.bashy"
Write-Host $bashyHome
New-Item -ItemType Directory -Force -Path $bashyHome\bin | Select-Object "bashy home created"
Copy-Item ./out/debug/bashy -Destination $bashyHome\bin\bashy.exe -Recurse



$old = [Environment]::GetEnvironmentVariable("Path", [System.EnvironmentVariableTarget]::User)
if ($old -like "*$bashyHome\bin*") { 
    Write-Host "bashy home already in PATH"
  }
  else{
    Write-Host "bashy home added to PATH"
    $new  =  "$old;$bashyHome\bin"
    [Environment]::SetEnvironmentVariable("Path",  $new, [System.EnvironmentVariableTarget]::User)
  }


  $old = [Environment]::GetEnvironmentVariable("Path", [System.EnvironmentVariableTarget]::User)
  if ($old -like "*$bashyHome\bashy.exe*") { 
    Write-Host "bashy home already in PATH"
  }
  else{
    Write-Host "bashy home added to PATH"
    $new  =  "$old;$bashyHome\bashy.exe"
    [Environment]::SetEnvironmentVariable("Path",  $new, [System.EnvironmentVariableTarget]::User)
  }
