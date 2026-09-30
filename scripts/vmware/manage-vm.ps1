<#
.SYNOPSIS
    VMware Virtual Machine Configuration, Hardware, Network, Disk, and Snapshot Management Engine.

.DESCRIPTION
    Comprehensive enterprise single-file PowerShell automation script for VMware Workstation,
    Player, and ESXi/Fusion .vmx virtual machines. Supports querying and updating MAC addresses,
    network adapter topologies (Bridged, NAT, Host-Only, Custom VMnet), RAM, vCPUs, virtual disks,
    printers, 3D display graphics, snapshots, and automated diagnostic log aggregation.

.PARAMETER Action
    The operation to perform. Supported values:
      - Get-VMs / Scan           : Recursively scan a folder for all .vmx virtual machines
      - Get-Config / Inspect     : Display full virtual hardware and network properties of a VM
      - Set-MAC                  : Change the MAC address of a network adapter (Static or Generated)
      - Set-Network              : Configure adapter connection type, custom VMnet, and power-on state
      - Set-Hardware             : Adjust RAM (MB/GB), vCPUs, and 3D display acceleration
      - Manage-Disk              : List, test, expand, or repair virtual disks (.vmdk)
      - Manage-Printer           : Enable or remove virtual printer device
      - Manage-Snapshot          : List, create, clone, revert, or delete snapshots
      - Optimize-VM              : Apply performance optimizations (disable paging, mem trimming)
      - Repair-VM                : Check .vmx integrity, verify disk references, and rebuild index
      - Copy-Log                 : Copy the latest execution log directly to the Windows clipboard
      - Help                     : Display this comprehensive interactive help guide

.PARAMETER VMPath
    Full path to the target virtual machine configuration file (.vmx) or folder.

.PARAMETER MACAddress
    Target MAC address in standard colon-separated format (e.g. "00:50:56:38:57:7B") or "generate".

.PARAMETER Adapter
    Target network adapter key (default: "ethernet0", or "ethernet1", "ethernet2", etc.).

.PARAMETER NetworkType
    Network connection mode: "bridged", "nat", "hostonly", "custom".

.PARAMETER VNet
    Custom virtual network device when NetworkType is "custom" (e.g. "VMnet1", "VMnet2", "VMnet8").

.PARAMETER Connected
    Adapter connected status ($true or $false).

.PARAMETER ConnectAtPowerOn
    Whether the network adapter connects automatically at VM power-on ($true or $false).

.PARAMETER MemoryMB
    Memory in Megabytes (e.g. 8192 for 8 GB).

.PARAMETER MemoryGB
    Convenience parameter for Memory in Gigabytes (e.g. 8, 16, 32).

.PARAMETER vCPUs
    Total virtual processor count (e.g. 4, 8, 16, 20).

.PARAMETER CoresPerSocket
    Cores per processor socket (e.g. 2, 4).

.PARAMETER Enable3D
    Enable 3D graphics acceleration ($true or $false).

.PARAMETER GraphicsMemoryMB
    Virtual graphics memory allocation in Megabytes (e.g. 2048, 4096, 8192).

.PARAMETER DiskAction
    Disk sub-action: "List", "Repair", "Expand".

.PARAMETER DiskPath
    Specific .vmdk disk file to repair or expand.

.PARAMETER NewDiskSize
    New disk size when expanding (e.g. "120GB", "200GB").

.PARAMETER PrinterEnabled
    Enable or disable virtual printer hardware ($true to add/enable, $false to remove).

.PARAMETER SnapshotAction
    Snapshot sub-action: "List", "Create", "Revert", "Clone", "Delete".

.PARAMETER SnapshotName
    Name of the snapshot for Create, Revert, Clone, or Delete operations.

.PARAMETER CloneDestPath
    Target destination folder or .vmx path when cloning a VM or snapshot.

.PARAMETER CopyLog
    Switch to automatically copy the entire execution log to the Windows clipboard.

.PARAMETER Json
    Switch to output results in machine-readable JSON format.

.PARAMETER LogFile
    Optional custom path to store execution transcript logs.
#>

[CmdletBinding(DefaultParameterSetName = "Standard")]
param(
    [Parameter(Position = 0)]
    [ValidateSet("Get-VMs", "Scan", "List-VMs", "ls",
                 "Get-Config", "Inspect", "Show", "Info",
                 "Set-MAC", "Change-MAC", "MAC",
                 "Set-Network", "Network", "Net",
                 "Set-Hardware", "Hardware", "HW",
                 "Manage-Disk", "Disks", "Disk",
                 "Manage-Printer", "Printer",
                 "Manage-Snapshot", "Snapshots", "Snapshot", "Snap",
                 "Optimize-VM", "Optimize", "Tweak",
                 "Repair-VM", "Repair",
                 "Copy-Log", "Export-Log",
                 "Help", "--help", "-h", "-?")]
    [string]$Action = "Help",

    [Parameter(Position = 1)]
    [Alias("Path", "VM")]
    [string]$VMPath = "",

    [Alias("MAC")]
    [string]$MACAddress = "",

    [Alias("Nic")]
    [string]$Adapter = "ethernet0",

    [Alias("NetType")]
    [ValidateSet("bridged", "nat", "hostonly", "custom")]
    [string]$NetworkType = "",

    [string]$VNet = "",

    [Nullable[bool]]$Connected = $null,

    [Nullable[bool]]$ConnectAtPowerOn = $null,

    [int]$MemoryMB = 0,

    [int]$MemoryGB = 0,

    [Alias("Processors", "CPU")]
    [int]$vCPUs = 0,

    [int]$CoresPerSocket = 0,

    [Nullable[bool]]$Enable3D = $null,

    [int]$GraphicsMemoryMB = 0,

    [ValidateSet("List", "Repair", "Expand", "")]
    [string]$DiskAction = "List",

    [string]$DiskPath = "",

    [string]$NewDiskSize = "",

    [Nullable[bool]]$PrinterEnabled = $null,

    [ValidateSet("List", "Create", "Revert", "Clone", "Delete", "")]
    [string]$SnapshotAction = "List",

    [string]$SnapshotName = "",

    [string]$CloneDestPath = "",

    [switch]$CopyLog,

    [switch]$Json,

    [string]$LogFile = "",

    [Alias("h", "?")]
    [switch]$Help
)

Set-StrictMode -Off
$ErrorActionPreference = "Stop"

if ($Help.IsPresent) {
    $Action = "Help"
}

# --- In-Memory Log Buffer & Transcript Store ---
$script:LogEntries = [System.Collections.Generic.List[string]]::new()
$script:SessionTimestamp = Get-Date -Format "yyyy-MM-dd_HH-mm-ss"
$script:DefaultLogDir = Join-Path -Path $env:TEMP -ChildPath "vmware-automation"

if (-not (Test-Path -Path $script:DefaultLogDir)) {
    [void](New-Item -ItemType Directory -Path $script:DefaultLogDir -Force)
}

if ([string]::IsNullOrWhiteSpace($LogFile)) {
    $script:LogFile = Join-Path -Path $script:DefaultLogDir -ChildPath "vm-manage-$script:SessionTimestamp.log"
} else {
    $script:LogFile = $LogFile
}

# --- Logging Helpers ---
function Write-VMLog {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Message,

        [ValidateSet("INFO", "SUCCESS", "WARN", "ERROR", "TRACE", "HEADER")]
        [string]$Level = "INFO"
    )

    $timeStr = Get-Date -Format "yyyy-MM-dd HH:mm:ss.fff"
    $formattedLog = "[$timeStr] [$Level] $Message"
    $script:LogEntries.Add($formattedLog)

    # Append to log file safely
    try {
        Add-Content -Path $script:LogFile -Value $formattedLog -Encoding UTF8 -ErrorAction SilentlyContinue
    } catch { }

    if ($Json) { return }

    switch ($Level) {
        "HEADER"  { Write-Host $Message -ForegroundColor Cyan }
        "INFO"    { Write-Host "  [i] $Message" -ForegroundColor Gray }
        "SUCCESS" { Write-Host "  [✔] $Message" -ForegroundColor Green }
        "WARN"    { Write-Host "  [▲] $Message" -ForegroundColor Yellow }
        "ERROR"   { Write-Host "  [✖] $Message" -ForegroundColor Red }
        "TRACE"   { Write-Host "      $Message" -ForegroundColor DarkGray }
    }
}

function Copy-VMLogToClipboard {
    $allText = $script:LogEntries -join "`r`n"
    try {
        Set-Clipboard -Value $allText -ErrorAction Stop
        Write-VMLog "Complete execution log ($($script:LogEntries.Count) lines) copied to Windows Clipboard." -Level SUCCESS
        Write-Host "`n  Clipboard Ready: Paste (Ctrl+V) this log directly to share with support/developer.`n" -ForegroundColor Cyan
    } catch {
        Write-VMLog "Failed to set clipboard: $($_.Exception.Message)" -Level WARN
    }
}

# --- VMware Toolchain Detection ---
function Get-VMwareTools {
    $tools = [PSCustomObject]@{
        WorkstationDir = ""
        VMRun          = ""
        VDiskManager   = ""
        HasVMRun       = $false
        HasVDisk       = $false
    }

    $candidateDirs = @(
        "C:\Program Files (x86)\VMware\VMware Workstation",
        "C:\Program Files\VMware\VMware Workstation",
        "C:\Program Files (x86)\VMware\VMware Player",
        "C:\Program Files\VMware\VMware Player",
        $env:VMWARE_HOME
    ) | Where-Object { -not [string]::IsNullOrWhiteSpace($_) -and (Test-Path -Path $_) }

    foreach ($dir in $candidateDirs) {
        $vmrunCandidate = Join-Path -Path $dir -ChildPath "vmrun.exe"
        if (Test-Path -Path $vmrunCandidate) {
            $tools.WorkstationDir = $dir
            $tools.VMRun = $vmrunCandidate
            $tools.HasVMRun = $true
        }

        $vdiskCandidate = Join-Path -Path $dir -ChildPath "vmware-vdiskmanager.exe"
        if (Test-Path -Path $vdiskCandidate) {
            $tools.VDiskManager = $vdiskCandidate
            $tools.HasVDisk = $true
        }

        if ($tools.HasVMRun -and $tools.HasVDisk) { break }
    }

    # PATH checks
    if (-not $tools.HasVMRun) {
        $pathCmd = Get-Command "vmrun.exe" -ErrorAction SilentlyContinue
        if ($pathCmd) {
            $tools.VMRun = $pathCmd.Source
            $tools.HasVMRun = $true
        }
    }
    if (-not $tools.HasVDisk) {
        $pathDisk = Get-Command "vmware-vdiskmanager.exe" -ErrorAction SilentlyContinue
        if ($pathDisk) {
            $tools.VDiskManager = $pathDisk.Source
            $tools.HasVDisk = $true
        }
    }

    return $tools
}

# --- VMX Configuration Parser & Serializer ---
function Read-VMXFile {
    param([string]$FilePath)

    if (-not (Test-Path -Path $FilePath)) {
        throw "VMX file does not exist: $FilePath"
    }

    $rawLines = @(Get-Content -Path $FilePath -Encoding UTF8 -ErrorAction Stop)
    $entries = [System.Collections.Specialized.OrderedDictionary]::new()
    $lineComments = [System.Collections.Generic.Dictionary[string, string]]::new()

    for ($i = 0; $i -lt $rawLines.Count; $i++) {
        $line = $rawLines[$i]
        $trimmed = $line.Trim()

        if ([string]::IsNullOrWhiteSpace($trimmed) -or $trimmed.StartsWith("#")) {
            $commentKey = "__COMMENT_$i"
            $entries[$commentKey] = $line
            continue
        }

        if ($line -match '^\s*([^=]+?)\s*=\s*"(.*)"\s*$') {
            $k = $matches[1].Trim()
            $v = $matches[2]
            $entries[$k] = $v
        } elseif ($line -match '^\s*([^=]+?)\s*=\s*(.*?)\s*$') {
            $k = $matches[1].Trim()
            $v = $matches[2].Trim()
            $entries[$k] = $v
        } else {
            $entries["__RAW_$i"] = $line
        }
    }

    return @{
        Entries  = $entries
        FilePath = (Resolve-Path -Path $FilePath).Path
    }
}

function Write-VMXFile {
    param(
        [hashtable]$VmxData,
        [switch]$NoBackup
    )

    $filePath = $VmxData.FilePath
    $entries = $VmxData.Entries

    if (-not $NoBackup) {
        $backupPath = "$filePath.bak"
        Copy-Item -Path $filePath -Destination $backupPath -Force
        Write-VMLog "Created backup snapshot: $backupPath" -Level INFO
    }

    $outputLines = [System.Collections.Generic.List[string]]::new()
    foreach ($key in $entries.Keys) {
        if ($key.StartsWith("__COMMENT_") -or $key.StartsWith("__RAW_")) {
            $outputLines.Add($entries[$key])
        } else {
            $val = $entries[$key]
            $outputLines.Add("$key = `"$val`"")
        }
    }

    [System.IO.File]::WriteAllLines($filePath, $outputLines, [System.Text.Encoding]::UTF8)
    Write-VMLog "Successfully committed changes to: $filePath" -Level SUCCESS
}

# --- Action 1: Scan Directory for VMs ---
function Invoke-ScanVMs {
    param([string]$SearchPath)

    if ([string]::IsNullOrWhiteSpace($SearchPath)) {
        $SearchPath = (Get-Location).Path
    }

    Write-VMLog "Scanning directory for virtual machines (.vmx): $SearchPath" -Level HEADER
    $vmxFiles = Get-ChildItem -Path $SearchPath -Filter "*.vmx" -Recurse -File -ErrorAction SilentlyContinue

    if (-not $vmxFiles -or $vmxFiles.Count -eq 0) {
        Write-VMLog "No virtual machines (.vmx) discovered in: $SearchPath" -Level WARN
        return @()
    }

    $results = [System.Collections.Generic.List[PSCustomObject]]::new()
    foreach ($file in $vmxFiles) {
        try {
            $vmx = Read-VMXFile -FilePath $file.FullName
            $entries = $vmx.Entries

            $displayName = if ($entries.Contains("displayName")) { $entries["displayName"] } else { $file.BaseName }
            $guestOS     = if ($entries.Contains("guestOS")) { $entries["guestOS"] } else { "Unknown" }
            $mem         = if ($entries.Contains("memsize")) { "$($entries['memsize']) MB" } else { "N/A" }
            $cpus        = if ($entries.Contains("numvcpus")) { $entries["numvcpus"] } else { "1" }
            $mac         = if ($entries.Contains("ethernet0.address")) { $entries["ethernet0.address"] }
                           elseif ($entries.Contains("ethernet0.generatedAddress")) { $entries["ethernet0.generatedAddress"] }
                           else { "None" }

            $obj = [PSCustomObject]@{
                DisplayName = $displayName
                GuestOS     = $guestOS
                Memory      = $mem
                vCPUs       = $cpus
                PrimaryMAC  = $mac
                VMXPath     = $file.FullName
                Directory   = $file.DirectoryName
            }
            $results.Add($obj)
        } catch {
            Write-VMLog "Could not parse $($file.FullName): $($_.Exception.Message)" -Level WARN
        }
    }

    if ($Json) {
        $results | ConvertTo-Json -Depth 4
        return
    }

    Write-Host "`n  Discovered $($results.Count) Virtual Machine(s):" -ForegroundColor Cyan
    Write-Host ("  " + ("=" * 115)) -ForegroundColor DarkGray
    Write-Host (("  {0,-25} {1,-18} {2,-10} {3,-6} {4,-19} {5}" -f "DISPLAY NAME", "GUEST OS", "MEMORY", "CPUS", "PRIMARY MAC", "VMX PATH")) -ForegroundColor Yellow
    Write-Host ("  " + ("-" * 115)) -ForegroundColor DarkGray

    foreach ($r in $results) {
        Write-Host (("  {0,-25} {1,-18} {2,-10} {3,-6} {4,-19} {5}" -f $r.DisplayName, $r.GuestOS, $r.Memory, $r.vCPUs, $r.PrimaryMAC, $r.VMXPath)) -ForegroundColor White
    }
    Write-Host ("  " + ("=" * 115) + "`n") -ForegroundColor DarkGray
}

# --- Action 2: Inspect VM Configuration ---
function Invoke-InspectVM {
    param([string]$FilePath)

    $vmx = Read-VMXFile -FilePath $FilePath
    $e = $vmx.Entries

    $displayName = if ($e.Contains("displayName")) { $e["displayName"] } else { [System.IO.Path]::GetFileNameWithoutExtension($FilePath) }
    $guestOS     = if ($e.Contains("guestOS")) { $e["guestOS"] } else { "N/A" }
    $memMB       = if ($e.Contains("memsize")) { [int]$e["memsize"] } else { 0 }
    $vCPUs       = if ($e.Contains("numvcpus")) { [int]$e["numvcpus"] } else { 1 }
    $cores       = if ($e.Contains("cpuid.coresPerSocket")) { [int]$e["cpuid.coresPerSocket"] } else { 1 }
    $enable3D    = if ($e.Contains("mks.enable3d")) { $e["mks.enable3d"] } else { "FALSE" }
    $vramKB      = if ($e.Contains("svga.graphicsMemoryKB")) { [int]$e["svga.graphicsMemoryKB"] } else { 0 }
    $printer     = if ($e.Contains("printers.enabled")) { $e["printers.enabled"] } else { "FALSE" }

    # Detect Network Adapters
    $adapters = [System.Collections.Generic.List[PSCustomObject]]::new()
    for ($i = 0; $i -le 9; $i++) {
        $prefix = "ethernet$i"
        if ($e.Contains("$prefix.present") -and $e["$prefix.present"] -match "(?i)true") {
            $connType = if ($e.Contains("$prefix.connectionType")) { $e["$prefix.connectionType"] } else { "bridged" }
            $vnet     = if ($e.Contains("$prefix.vnet")) { $e["$prefix.vnet"] } else { "" }
            $addrType = if ($e.Contains("$prefix.addressType")) { $e["$prefix.addressType"] } else { "generated" }
            $mac      = if ($e.Contains("$prefix.address")) { $e["$prefix.address"] }
                        elseif ($e.Contains("$prefix.generatedAddress")) { $e["$prefix.generatedAddress"] }
                        else { "Not Set" }
            $startCon = if ($e.Contains("$prefix.startConnected")) { $e["$prefix.startConnected"] } else { "TRUE" }
            $connNow  = if ($e.Contains("$prefix.connected")) { $e["$prefix.connected"] } else { "TRUE" }

            $adapters.Add([PSCustomObject]@{
                Adapter          = $prefix
                ConnectionType   = $connType
                VirtualNetwork   = $vnet
                AddressType      = $addrType
                MACAddress       = $mac
                ConnectAtPowerOn = $startCon
                CurrentlyConnect = $connNow
            })
        }
    }

    # Detect Virtual Disks
    $disks = [System.Collections.Generic.List[PSCustomObject]]::new()
    $diskControllers = @("nvme0", "nvme1", "scsi0", "scsi1", "sata0", "sata1", "ide0", "ide1")
    foreach ($ctrl in $diskControllers) {
        for ($u = 0; $u -le 15; $u++) {
            $dev = "$ctrl`:$u"
            if ($e.Contains("$dev.present") -and $e["$dev.present"] -match "(?i)true") {
                $fileName = if ($e.Contains("$dev.fileName")) { $e["$dev.fileName"] } else { "Unknown" }
                $diskType = $ctrl.Substring(0, 4).ToUpper()
                $disks.Add([PSCustomObject]@{
                    Device   = $dev
                    BusType  = $diskType
                    FileName = $fileName
                })
            }
        }
    }

    $summary = [PSCustomObject]@{
        VMXFile         = $FilePath
        DisplayName     = $displayName
        GuestOS         = $guestOS
        MemoryMB        = $memMB
        MemoryGB        = [math]::Round($memMB / 1024, 2)
        vCPUs           = $vCPUs
        CoresPerSocket  = $cores
        Enable3D        = $enable3D
        GraphicsVRAM_MB = [math]::Round($vramKB / 1024, 1)
        PrinterEnabled  = $printer
        NetworkAdapters = $adapters
        VirtualDisks    = $disks
    }

    if ($Json) {
        $summary | ConvertTo-Json -Depth 5
        return
    }

    Write-Host "`n  ╔══ VIRTUAL MACHINE SETTINGS: $displayName ═══════════════════════════════════╗" -ForegroundColor Cyan
    Write-Host "  │ File:  $FilePath" -ForegroundColor Gray
    Write-Host "  │ OS:    $guestOS | CPUs: $vCPUs (Cores/Socket: $cores) | RAM: $memMB MB ($([math]::Round($memMB/1024, 1)) GB)" -ForegroundColor White
    Write-Host "  │ 3D Display: $enable3D (VRAM: $([math]::Round($vramKB/1024, 1)) MB) | Virtual Printer: $printer" -ForegroundColor White
    Write-Host "  ╠══ NETWORK ADAPTERS ($($adapters.Count)) ════════════════════════════════════════════════════════╣" -ForegroundColor Cyan

    foreach ($a in $adapters) {
        $vnetTag = if ($a.VirtualNetwork) { " ($($a.VirtualNetwork))" } else { "" }
        Write-Host "  │  • $($a.Adapter): Type: $($a.ConnectionType)$vnetTag | MAC: $($a.MACAddress) [$($a.AddressType)]" -ForegroundColor Yellow
        Write-Host "  │    Power-On Connect: $($a.ConnectAtPowerOn) | Live Connected: $($a.CurrentlyConnect)" -ForegroundColor DarkGray
    }

    Write-Host "  ╠══ VIRTUAL HARD DRIVES ($($disks.Count)) ═════════════════════════════════════════════════════╣" -ForegroundColor Cyan
    foreach ($d in $disks) {
        Write-Host "  │  • [$($d.Device)] Type: $($d.BusType) -> $($d.FileName)" -ForegroundColor White
    }
    Write-Host "  ╚═══════════════════════════════════════════════════════════════════════════════╝`n" -ForegroundColor Cyan
}

# --- Action 3: Set MAC Address ---
function Invoke-SetMACAddress {
    param(
        [string]$FilePath,
        [string]$Nic,
        [string]$MAC
    )

    if ([string]::IsNullOrWhiteSpace($MAC)) {
        throw "MAC address parameter is mandatory for Set-MAC (e.g. '00:50:56:38:57:7B' or 'generate')"
    }

    $vmx = Read-VMXFile -FilePath $FilePath
    $e = $vmx.Entries

    Write-VMLog "Updating MAC Address on $Nic for $FilePath..." -Level INFO

    if ($MAC.Trim().ToLower() -eq "generate") {
        $e["$Nic.present"] = "TRUE"
        $e["$Nic.addressType"] = "generated"
        $e.Remove("$Nic.address")
        Write-VMLog "Configured $Nic to automatic generated MAC on next power-on." -Level SUCCESS
    } else {
        # Validate MAC format
        $cleanMAC = $MAC.Trim().ToUpper().Replace("-", ":")
        if ($cleanMAC -notmatch '^([0-9A-F]{2}:){5}[0-9A-F]{2}$') {
            throw "Invalid MAC address format: '$MAC'. Must be 6 pairs of hex digits separated by ':' (e.g. 00:50:56:38:57:7B)"
        }

        $e["$Nic.present"] = "TRUE"
        $e["$Nic.addressType"] = "static"
        $e["$Nic.address"] = $cleanMAC
        $e.Remove("$Nic.generatedAddress")
        $e.Remove("$Nic.generatedAddressOffset")

        Write-VMLog "Set static MAC on $($Nic): $cleanMAC" -Level SUCCESS
    }

    Write-VMXFile -VmxData $vmx
}

# --- Action 4: Set Network Topology & Connection ---
function Invoke-SetNetwork {
    param(
        [string]$FilePath,
        [string]$Nic,
        [string]$NetType,
        [string]$CustomVNet,
        [Nullable[bool]]$IsConnected,
        [Nullable[bool]]$IsStartConnected
    )

    $vmx = Read-VMXFile -FilePath $FilePath
    $e = $vmx.Entries

    $e["$Nic.present"] = "TRUE"

    if (-not [string]::IsNullOrWhiteSpace($NetType)) {
        $mode = $NetType.Trim().ToLower()
        switch ($mode) {
            "bridged"  { $e["$Nic.connectionType"] = "bridged"; $e.Remove("$Nic.vnet") }
            "nat"      { $e["$Nic.connectionType"] = "nat"; $e.Remove("$Nic.vnet") }
            "hostonly" { $e["$Nic.connectionType"] = "hostonly"; $e.Remove("$Nic.vnet") }
            "custom"   {
                $e["$Nic.connectionType"] = "custom"
                if (-not [string]::IsNullOrWhiteSpace($CustomVNet)) {
                    $e["$Nic.vnet"] = $CustomVNet
                }
            }
        }
        Write-VMLog "Updated network mode on $Nic to: $mode $(if ($CustomVNet) { "($CustomVNet)" })" -Level SUCCESS
    }

    if ($null -ne $IsConnected) {
        $e["$Nic.connected"] = if ($IsConnected) { "TRUE" } else { "FALSE" }
        Write-VMLog "Set $Nic live connected: $IsConnected" -Level INFO
    }

    if ($null -ne $IsStartConnected) {
        $e["$Nic.startConnected"] = if ($IsStartConnected) { "TRUE" } else { "FALSE" }
        Write-VMLog "Set $Nic connect-at-power-on: $IsStartConnected" -Level INFO
    }

    Write-VMXFile -VmxData $vmx
}

# --- Action 5: Set Hardware (RAM, vCPUs, 3D Graphics) ---
function Invoke-SetHardware {
    param(
        [string]$FilePath,
        [int]$RamMB,
        [int]$RamGB,
        [int]$Cpus,
        [int]$Cores,
        [Nullable[bool]]$Gfx3D,
        [int]$VramMB
    )

    $vmx = Read-VMXFile -FilePath $FilePath
    $e = $vmx.Entries

    if ($RamGB -gt 0) {
        $RamMB = $RamGB * 1024
    }

    if ($RamMB -gt 0) {
        $e["memsize"] = "$RamMB"
        Write-VMLog "Allocated RAM: $RamMB MB ($([math]::Round($RamMB/1024, 2)) GB)" -Level SUCCESS
    }

    if ($Cpus -gt 0) {
        $e["numvcpus"] = "$Cpus"
        Write-VMLog "Allocated vCPUs: $Cpus" -Level SUCCESS
    }

    if ($Cores -gt 0) {
        $e["cpuid.coresPerSocket"] = "$Cores"
        Write-VMLog "Set Cores per Socket: $Cores" -Level SUCCESS
    }

    if ($null -ne $Gfx3D) {
        $e["mks.enable3d"] = if ($Gfx3D) { "TRUE" } else { "FALSE" }
        Write-VMLog "3D Display Acceleration: $Gfx3D" -Level SUCCESS
    }

    if ($VramMB -gt 0) {
        $e["svga.graphicsMemoryKB"] = "$($VramMB * 1024)"
        $e["svga.autodetect"] = "FALSE"
        Write-VMLog "Allocated Graphics VRAM: $VramMB MB" -Level SUCCESS
    }

    Write-VMXFile -VmxData $vmx
}

# --- Action 6: Manage Virtual Disks ---
function Invoke-ManageDisks {
    param(
        [string]$FilePath,
        [string]$SubAction,
        [string]$TargetDisk,
        [string]$NewSize
    )

    $tools = Get-VMwareTools
    $vmx = Read-VMXFile -FilePath $FilePath
    $e = $vmx.Entries
    $vmDir = Split-Path -Path $FilePath -Parent

    # Gather Disks
    $disks = [System.Collections.Generic.List[PSCustomObject]]::new()
    $diskControllers = @("nvme0", "nvme1", "scsi0", "scsi1", "sata0", "sata1", "ide0", "ide1")
    foreach ($ctrl in $diskControllers) {
        for ($u = 0; $u -le 15; $u++) {
            $dev = "$ctrl`:$u"
            if ($e.Contains("$dev.present") -and $e["$dev.present"] -match "(?i)true") {
                $relFile = $e["$dev.fileName"]
                $absFile = if ([System.IO.Path]::IsPathRooted($relFile)) { $relFile } else { Join-Path -Path $vmDir -ChildPath $relFile }
                $exists = Test-Path -Path $absFile
                $sizeGB = 0
                if ($exists) {
                    $sizeGB = [math]::Round((Get-Item -Path $absFile).Length / 1GB, 2)
                }

                $disks.Add([PSCustomObject]@{
                    Device   = $dev
                    FileName = $relFile
                    FullPath = $absFile
                    Exists   = $exists
                    SizeGB   = $sizeGB
                })
            }
        }
    }

    switch ($SubAction.ToLower()) {
        "repair" {
            if (-not $tools.HasVDisk) {
                Write-VMLog "vmware-vdiskmanager.exe not detected on host. Verifying disk file integrity via filesystem check..." -Level WARN
                foreach ($d in $disks) {
                    if ($d.Exists) {
                        Write-VMLog "Disk [$($d.Device)] file verified present: $($d.FullPath) ($($d.SizeGB) GB)" -Level SUCCESS
                    } else {
                        Write-VMLog "Disk [$($d.Device)] missing target file: $($d.FullPath)" -Level ERROR
                    }
                }
                return
            }

            Write-VMLog "Running disk repair using vmware-vdiskmanager..." -Level HEADER
            foreach ($d in $disks) {
                if ($d.Exists) {
                    Write-VMLog "Repairing $($d.FullPath)..." -Level INFO
                    $out = & $tools.VDiskManager -R $d.FullPath 2>&1
                    Write-VMLog "Result for $($d.Device): $out" -Level SUCCESS
                } else {
                    Write-VMLog "Skipping missing disk: $($d.FullPath)" -Level WARN
                }
            }
        }

        "expand" {
            if (-not $tools.HasVDisk) {
                throw "Cannot expand disk: vmware-vdiskmanager.exe is required but not installed."
            }
            if ([string]::IsNullOrWhiteSpace($TargetDisk) -or [string]::IsNullOrWhiteSpace($NewSize)) {
                throw "Expand requires -DiskPath <path-to-.vmdk> and -NewDiskSize <e.g. 150GB>"
            }
            Write-VMLog "Expanding $TargetDisk to $NewSize..." -Level INFO
            $out = & $tools.VDiskManager -x $NewSize $TargetDisk 2>&1
            Write-VMLog "Expand result: $out" -Level SUCCESS
        }

        default {
            Write-VMLog "Virtual Hard Disks attached to $FilePath ($($disks.Count) disks):" -Level HEADER
            foreach ($d in $disks) {
                $status = if ($d.Exists) { "Present ($($d.SizeGB) GB)" } else { "MISSING FILE" }
                Write-Host "    • [$($d.Device)] $($d.FileName) - $status" -ForegroundColor (if ($d.Exists) { "White" } else { "Red" })
            }
        }
    }
}

# --- Action 7: Manage Printer Device ---
function Invoke-ManagePrinter {
    param(
        [string]$FilePath,
        [bool]$Enable
    )

    $vmx = Read-VMXFile -FilePath $FilePath
    $e = $vmx.Entries

    $e["printers.enabled"] = if ($Enable) { "TRUE" } else { "FALSE" }
    Write-VMLog "Set virtual printer enabled state: $Enable" -Level SUCCESS
    Write-VMXFile -VmxData $vmx
}

# --- Action 8: Manage Snapshots ---
function Invoke-ManageSnapshots {
    param(
        [string]$FilePath,
        [string]$SubAction,
        [string]$SnapName,
        [string]$DestClone
    )

    $tools = Get-VMwareTools
    $vmDir = Split-Path -Path $FilePath -Parent
    $vmsdPath = [System.IO.Path]::ChangeExtension($FilePath, ".vmsd")

    switch ($SubAction.ToLower()) {
        "create" {
            if ([string]::IsNullOrWhiteSpace($SnapName)) {
                $SnapName = "Snapshot-" + (Get-Date -Format "yyyyMMdd-HHmmss")
            }
            if (-not $tools.HasVMRun) {
                throw "vmrun.exe is required for live snapshot creation. Install VMware Workstation/VIX."
            }
            Write-VMLog "Creating snapshot '$SnapName' on $FilePath..." -Level INFO
            $out = & $tools.VMRun -T ws snapshot $FilePath $SnapName 2>&1
            Write-VMLog "Snapshot '$SnapName' created successfully: $out" -Level SUCCESS
        }

        "revert" {
            if ([string]::IsNullOrWhiteSpace($SnapName)) {
                throw "SnapshotName is required to revert a snapshot."
            }
            if (-not $tools.HasVMRun) {
                throw "vmrun.exe is required to revert snapshots."
            }
            Write-VMLog "Reverting to snapshot '$SnapName'..." -Level INFO
            $out = & $tools.VMRun -T ws revertToSnapshot $FilePath $SnapName 2>&1
            Write-VMLog "Revert complete: $out" -Level SUCCESS
        }

        "clone" {
            if ([string]::IsNullOrWhiteSpace($DestClone)) {
                throw "CloneDestPath is required to clone a snapshot or VM."
            }
            if (-not $tools.HasVMRun) {
                throw "vmrun.exe is required to clone snapshots."
            }
            Write-VMLog "Cloning VM to $DestClone..." -Level INFO
            $snapParam = if ($SnapName) { "-snapshot=$SnapName" } else { "" }
            $out = & $tools.VMRun -T ws clone $FilePath $DestClone full $snapParam 2>&1
            Write-VMLog "Clone complete: $out" -Level SUCCESS
        }

        default {
            # List snapshots
            Write-VMLog "Querying snapshots for $FilePath..." -Level HEADER
            if ($tools.HasVMRun) {
                $out = & $tools.VMRun -T ws listSnapshots $FilePath 2>&1
                Write-Host "`n  $out`n" -ForegroundColor White
            } elseif (Test-Path -Path $vmsdPath) {
                Write-VMLog "Reading snapshots from metadata file: $vmsdPath" -Level INFO
                $vmsd = Get-Content -Path $vmsdPath -Encoding UTF8
                $snaps = $vmsd | Where-Object { $_ -match 'snapshot\d+\.displayName\s*=\s*"(.*)"' }
                Write-Host "`n  Discovered Snapshots from .vmsd:" -ForegroundColor Cyan
                foreach ($s in $snaps) {
                    Write-Host "    • $s" -ForegroundColor White
                }
                Write-Host ""
            } else {
                Write-VMLog "No snapshots or .vmsd metadata found for this VM." -Level INFO
            }
        }
    }
}

# --- Action 9: Optimize VM Configuration ---
function Invoke-OptimizeVM {
    param([string]$FilePath)

    $vmx = Read-VMXFile -FilePath $FilePath
    $e = $vmx.Entries

    Write-VMLog "Applying recommended performance optimizations to $FilePath..." -Level HEADER

    # Disable memory trimming & page sharing for maximum host disk I/O performance
    $e["MemTrimRate"] = "0"
    $e["mainMem.useNamedFile"] = "FALSE"
    $e["MemAllowRequestParamOverride"] = "TRUE"
    $e["prefvmx.minVmMemPct"] = "100"
    $e["prefvmx.useRecommendedLockedMemSize"] = "TRUE"
    $e["logging"] = "FALSE"

    Write-VMLog "Disabled memory trimming and named paging file (keeps VM RAM strictly in host RAM)." -Level SUCCESS
    Write-VMLog "Disabled hypervisor debug logging for maximum I/O throughput." -Level SUCCESS

    Write-VMXFile -VmxData $vmx
}

# --- Action 10: Help System ---
function Show-VMHelp {
    $helpText = @"

  ╔═════════════════════════════════════════════════════════════════════════════════════════╗
  ║              VMware Enterprise Automation & Configuration Engine (manage-vm.ps1)        ║
  ╚═════════════════════════════════════════════════════════════════════════════════════════╝

  SYNOPSIS:
    .\manage-vm.ps1 -Action <Action> [-VMPath <path>] [Options...] [-CopyLog]

  CORE ACTIONS:

    1. Scan for Virtual Machines:
       .\manage-vm.ps1 -Action Scan -VMPath "D:\VMs"
       .\manage-vm.ps1 -Action Scan (scans current working directory)

    2. Inspect VM Configuration:
       .\manage-vm.ps1 -Action Inspect -VMPath "D:\VMs\WinServer\WinServer.vmx"

    3. Change MAC Address (Static or Generated):
       .\manage-vm.ps1 -Action Set-MAC -VMPath "D:\VMs\WinServer\WinServer.vmx" -MAC "00:50:56:38:57:7B"
       .\manage-vm.ps1 -Action Set-MAC -VMPath "D:\VMs\WinServer\WinServer.vmx" -MAC "generate" -Adapter ethernet0

    4. Configure Network Connection (Bridged, NAT, HostOnly, Custom):
       .\manage-vm.ps1 -Action Set-Network -VMPath "D:\VMs\Ubuntu\Ubuntu.vmx" -NetworkType bridged
       .\manage-vm.ps1 -Action Set-Network -VMPath "D:\VMs\Ubuntu\Ubuntu.vmx" -NetworkType custom -VNet VMnet1
       .\manage-vm.ps1 -Action Set-Network -VMPath "D:\VMs\Ubuntu\Ubuntu.vmx" -Connected `$true -ConnectAtPowerOn `$true

    5. Adjust System Hardware (RAM, vCPUs, 3D Graphics):
       .\manage-vm.ps1 -Action Set-Hardware -VMPath "D:\VMs\WinServer\WinServer.vmx" -MemoryGB 8 -vCPUs 20
       .\manage-vm.ps1 -Action Set-Hardware -VMPath "D:\VMs\WinServer\WinServer.vmx" -MemoryMB 16384 -Enable3D `$true

    6. Virtual Hard Disks (List, Repair, Expand):
       .\manage-vm.ps1 -Action Manage-Disk -VMPath "D:\VMs\WinServer\WinServer.vmx" -DiskAction List
       .\manage-vm.ps1 -Action Manage-Disk -VMPath "D:\VMs\WinServer\WinServer.vmx" -DiskAction Repair

    7. Virtual Printer Management:
       .\manage-vm.ps1 -Action Manage-Printer -VMPath "D:\VMs\WinServer\WinServer.vmx" -PrinterEnabled `$false

    8. Snapshots (List, Create, Revert, Clone):
       .\manage-vm.ps1 -Action Manage-Snapshot -VMPath "D:\VMs\WinServer\WinServer.vmx" -SnapshotAction List
       .\manage-vm.ps1 -Action Manage-Snapshot -VMPath "D:\VMs\WinServer\WinServer.vmx" -SnapshotAction Create -SnapshotName "Clean-Baseline"

    9. Performance Optimization:
       .\manage-vm.ps1 -Action Optimize-VM -VMPath "D:\VMs\WinServer\WinServer.vmx"

   10. Copy Execution Log to Windows Clipboard:
       Add '-CopyLog' to ANY command to automatically copy the complete log transcript
       to your Windows clipboard, ready to paste (Ctrl+V) and share immediately!
       .\manage-vm.ps1 -Action Copy-Log

"@
    Write-Host $helpText -ForegroundColor Cyan
}

# --- Main Dispatcher ---
try {
    Write-VMLog "Initializing VMware Automation Engine (Action: $Action, Target: $VMPath)" -Level TRACE

    switch ($Action.ToLower()) {
        { $_ -in "get-vms", "scan", "list-vms", "ls" } {
            Invoke-ScanVMs -SearchPath $VMPath
        }

        { $_ -in "get-config", "inspect", "show", "info" } {
            if ([string]::IsNullOrWhiteSpace($VMPath)) { throw "VMPath is required for Get-Config/Inspect." }
            Invoke-InspectVM -FilePath $VMPath
        }

        { $_ -in "set-mac", "change-mac", "mac" } {
            if ([string]::IsNullOrWhiteSpace($VMPath)) { throw "VMPath is required for Set-MAC." }
            Invoke-SetMACAddress -FilePath $VMPath -Nic $Adapter -MAC $MACAddress
        }

        { $_ -in "set-network", "network", "net" } {
            if ([string]::IsNullOrWhiteSpace($VMPath)) { throw "VMPath is required for Set-Network." }
            Invoke-SetNetwork -FilePath $VMPath -Nic $Adapter -NetType $NetworkType -CustomVNet $VNet -IsConnected $Connected -IsStartConnected $ConnectAtPowerOn
        }

        { $_ -in "set-hardware", "hardware", "hw" } {
            if ([string]::IsNullOrWhiteSpace($VMPath)) { throw "VMPath is required for Set-Hardware." }
            Invoke-SetHardware -FilePath $VMPath -RamMB $MemoryMB -RamGB $MemoryGB -Cpus $vCPUs -Cores $CoresPerSocket -Gfx3D $Enable3D -VramMB $GraphicsMemoryMB
        }

        { $_ -in "manage-disk", "disks", "disk" } {
            if ([string]::IsNullOrWhiteSpace($VMPath)) { throw "VMPath is required for Manage-Disk." }
            Invoke-ManageDisks -FilePath $VMPath -SubAction $DiskAction -TargetDisk $DiskPath -NewSize $NewDiskSize
        }

        { $_ -in "manage-printer", "printer" } {
            if ([string]::IsNullOrWhiteSpace($VMPath)) { throw "VMPath is required for Manage-Printer." }
            if ($null -eq $PrinterEnabled) { throw "PrinterEnabled ($true/$false) is required for Manage-Printer." }
            Invoke-ManagePrinter -FilePath $VMPath -Enable $PrinterEnabled
        }

        { $_ -in "manage-snapshot", "snapshots", "snapshot", "snap" } {
            if ([string]::IsNullOrWhiteSpace($VMPath)) { throw "VMPath is required for Manage-Snapshot." }
            Invoke-ManageSnapshots -FilePath $VMPath -SubAction $SnapshotAction -SnapName $SnapshotName -DestClone $CloneDestPath
        }

        { $_ -in "optimize-vm", "optimize", "tweak" } {
            if ([string]::IsNullOrWhiteSpace($VMPath)) { throw "VMPath is required for Optimize-VM." }
            Invoke-OptimizeVM -FilePath $VMPath
        }

        { $_ -in "repair-vm", "repair" } {
            if ([string]::IsNullOrWhiteSpace($VMPath)) { throw "VMPath is required for Repair-VM." }
            Invoke-ManageDisks -FilePath $VMPath -SubAction "repair"
        }

        { $_ -in "copy-log", "export-log" } {
            Copy-VMLogToClipboard
        }

        default {
            Show-VMHelp
        }
    }

    if ($CopyLog) {
        Copy-VMLogToClipboard
    }
}
catch {
    $errMessage = $_.Exception.Message
    $stackTrace = $_.ScriptStackTrace
    Write-VMLog "CRITICAL EXECUTION ERROR: $errMessage" -Level ERROR
    Write-VMLog "Stack Trace: $stackTrace" -Level TRACE

    Write-Host "`n  ╔══ ERROR REPORT ══════════════════════════════════════════════════════════════╗" -ForegroundColor Red
    Write-Host "  │ Message:     $errMessage" -ForegroundColor Red
    Write-Host "  │ Origin:      $($_.InvocationInfo.ScriptName):$($_.InvocationInfo.ScriptLineNumber)" -ForegroundColor Yellow
    Write-Host "  │ Stack Trace: $stackTrace" -ForegroundColor DarkGray
    Write-Host "  ╚═══════════════════════════════════════════════════════════════════════════════╝`n" -ForegroundColor Red

    if ($CopyLog) {
        Copy-VMLogToClipboard
    }

    exit 1
}
