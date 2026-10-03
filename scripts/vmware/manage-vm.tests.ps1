# Pester unit and integration tests for manage-vm.ps1
# Requires Pester 3.4.0+ / PowerShell Core (pwsh)

Describe "VMware Automation Engine - manage-vm.ps1" {

    BeforeAll {
        $script:scriptPath = Join-Path -Path $PSScriptRoot -ChildPath "manage-vm.ps1"
        $script:testDir = Join-Path -Path ([System.IO.Path]::GetTempPath()) -ChildPath "pester_vmware_test_$([System.Guid]::NewGuid().ToString('N'))"
        New-Item -ItemType Directory -Path $script:testDir -Force | Out-Null

        $script:vmxFile = Join-Path -Path $script:testDir -ChildPath "test-machine.vmx"
        $vmxContent = @"
.encoding = "windows-1252"
config.version = "8"
virtualHW.version = "19"
displayName = "Pester-Test-VM"
guestOS = "windows10-64"
memsize = "4096"
numvcpus = "2"
ethernet0.present = "TRUE"
ethernet0.connectionType = "nat"
ethernet0.addressType = "generated"
ethernet0.generatedAddress = "00:0c:29:11:22:33"
"@
        Set-Content -Path $script:vmxFile -Value $vmxContent -Encoding ASCII
    }

    AfterAll {
        if (Test-Path $script:testDir) {
            Remove-Item -Recurse -Force $script:testDir -ErrorAction SilentlyContinue
        }
    }

    Context "Script Integrity and Help System" {

        It "manage-vm.ps1 script file exists on disk" {
            Test-Path $script:scriptPath | Should Be $true
        }

        It "parses with zero PowerShell syntax errors" {
            $parseErrors = $null
            $tokens = $null
            [System.Management.Automation.Language.Parser]::ParseFile($script:scriptPath, [ref]$tokens, [ref]$parseErrors) | Out-Null
            $parseErrors.Count | Should Be 0
        }

        It "renders interactive help banner with -Action Help" {
            $output = & pwsh -NoProfile -File $script:scriptPath -Action Help
            $output -join "`n" | Should Match "VMware Enterprise Automation"
        }

        It "validates parameter actions and rejects unknown action ListGroups" {
            $output = & pwsh -NoProfile -File $script:scriptPath -Action ListGroups 2>&1
            $isRejected = ($LASTEXITCODE -ne 0)
            $isRejected | Should Be $true
            $output -join "`n" | Should Match "Cannot validate argument on parameter 'Action'"
        }
    }

    Context "VM Configuration Inspection" {

        It "inspects mock VMX configuration correctly" {
            $output = & pwsh -NoProfile -File $script:scriptPath -Action Inspect -VMPath $script:vmxFile
            $joined = $output -join "`n"
            $joined | Should Match "Pester-Test-VM"
            $joined | Should Match "4096 MB"
        }

        It "discovers virtual machines in scan directory" {
            $output = & pwsh -NoProfile -File $script:scriptPath -Action Scan -VMPath $script:testDir
            $joined = $output -join "`n"
            $joined | Should Match "Discovered 1 Virtual Machine"
            $joined | Should Match "Pester-Test-VM"
        }
    }

    Context "MAC Mutation and Network Logic" {

        It "sets static MAC address and creates automatic safety backup" {
            $targetMAC = "00:50:56:38:57:7B"
            $output = & pwsh -NoProfile -File $script:scriptPath -Action Set-MAC -VMPath $script:vmxFile -MAC $targetMAC -Adapter ethernet0
            $joined = $output -join "`n"
            $joined | Should Match "Set static MAC on ethernet0"

            $bakFile = "$($script:vmxFile).bak"
            Test-Path $bakFile | Should Be $true

            $content = Get-Content $script:vmxFile -Raw
            $content | Should Match $targetMAC
        }

        It "configures dynamic generated MAC address" {
            $output = & pwsh -NoProfile -File $script:scriptPath -Action Set-MAC -VMPath $script:vmxFile -MAC "generate" -Adapter ethernet0
            $joined = $output -join "`n"
            $joined | Should Match "Configured ethernet0 to automatic generated MAC"

            $content = Get-Content $script:vmxFile -Raw
            $content | Should Match 'ethernet0.addressType = "generated"'
        }

        It "updates network connection type to bridged" {
            $output = & pwsh -NoProfile -File $script:scriptPath -Action Set-Network -VMPath $script:vmxFile -NetworkType bridged -Adapter ethernet0
            $joined = $output -join "`n"
            $joined | Should Match "Updated network mode on ethernet0 to: bridged"

            $content = Get-Content $script:vmxFile -Raw
            $content | Should Match 'ethernet0.connectionType = "bridged"'
        }
    }

    Context "Hardware Modification and Optimization" {

        It "adjusts RAM and vCPUs hardware parameters" {
            $output = & pwsh -NoProfile -File $script:scriptPath -Action Set-Hardware -VMPath $script:vmxFile -MemoryGB 8 -vCPUs 4
            $joined = $output -join "`n"
            $joined | Should Match "Allocated RAM: 8192 MB"
            $joined | Should Match "Allocated vCPUs: 4"

            $content = Get-Content $script:vmxFile -Raw
            $content | Should Match 'memsize = "8192"'
            $content | Should Match 'numvcpus = "4"'
        }

        It "applies performance tuning optimizations" {
            $output = & pwsh -NoProfile -File $script:scriptPath -Action Optimize-VM -VMPath $script:vmxFile
            $joined = $output -join "`n"
            $joined | Should Match "Disabled memory trimming"

            $content = Get-Content $script:vmxFile -Raw
            $content | Should Match 'mainMem.useNamedFile = "FALSE"'
        }
    }
}
