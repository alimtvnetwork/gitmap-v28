# gitmap os tweak

Configure Windows desktop tweaks, classic context menus, Start menu layouts, power schemes, and hibernation.

## Usage

```bash
gitmap os tweak [category] [action]
```

## Categories & Actions

| Category | Action | Description |
|----------|--------|-------------|
| status (st) | *(none)* | Inspect active system tweaks, power schemes, and hibernation status |
| context-menu | classic | Restore classic Windows 10 right-click context menu (bypasses Windows 11 XAML menu) |
| context-menu | modern | Restore modern default Windows 11 right-click context menu |
| start-menu | classic | Restore classic Start Menu layout |
| start-menu | default | Restore default Start Menu layout |
| power | ultimate | Duplicate and activate Windows Ultimate Performance power scheme |
| power | balanced | Activate default Windows Balanced power scheme |
| hibernate | off | Disable hibernation and delete C:\hiberfil.sys to reclaim disk space |
| hibernate | on | Enable hibernation |

## Examples

### Inspect Active Tweaks

```bash
gitmap os tweak status
```

### Restore Windows 10 Classic Right-Click Menu

```bash
gitmap os tweak context-menu classic
```

### Restore Windows 11 Modern Context Menu

```bash
gitmap os tweak context-menu modern
```

### Activate Ultimate Performance Power Scheme

```bash
gitmap os tweak power ultimate
```

### Turn Off Hibernation to Free RAM-Sized Disk Space

```bash
gitmap os tweak hibernate off
```
