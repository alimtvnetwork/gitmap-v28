# Pester unit tests for manage-vm.ps1
Describe "manage-vm.ps1" {
    It "script file exists" {
        Test-Path "$PSScriptRoot\manage-vm.ps1" | Should -Be $true
    }
}
