# ActionType — Reference Implementation

**Version:** 2.0.0
**Updated:** 2026-04-09

> **Purpose:** Second reference implementation using `match`-based metadata. Demonstrates a large enum (40+ cases) with per-case and group helpers.

---

## ActionType Enum

```php
namespace PluginName\Enums;

if (!defined('ABSPATH')) {
    exit;
}

enum ActionType: string
{
    // ── Core Plugin Actions ─────────────────────────────────
    case Upload           = 'Upload';
    case UploadActive     = 'UploadActive';
    case UploadInitiated  = 'UploadInitiated';
    case Enable           = 'Enable';
    case Disable          = 'Disable';
    case Delete           = 'Delete';
    case FileReplace      = 'FileReplace';
    case FileDelete       = 'FileDelete';
    case Sync             = 'Sync';
    case SyncDelete       = 'SyncDelete';

    // ── Post/Content Actions ────────────────────────────────
    case PostCreate       = 'PostCreate';
    case PostUpdate       = 'PostUpdate';
    case CategoryCreate   = 'CategoryCreate';
    case MediaUpload      = 'MediaUpload';

    // ── Auth ────────────────────────────────────────────────
    case AuthFailed       = 'AuthFailed';

    // ── Export Actions ──────────────────────────────────────
    case ExportSelf       = 'ExportSelf';
    case ExportPlugin     = 'ExportPlugin';

    // ── Plugin Backup Actions ───────────────────────────────
    case PluginBackup        = 'PluginBackup';
    case PluginBackupRestore = 'PluginBackupRestore';
    case PluginBackupDelete  = 'PluginBackupDelete';

    // ── Agent Actions ───────────────────────────────────────
    case AgentAdd           = 'AgentAdd';
    case AgentRemove        = 'AgentRemove';
    case AgentTest          = 'AgentTest';
    case AgentSync          = 'AgentSync';
    case AgentApiError      = 'AgentApiError';

    // ── Snapshot Actions ────────────────────────────────────
    case SnapshotCreate          = 'SnapshotCreate';
    case SnapshotRestore         = 'SnapshotRestore';
    case SnapshotDelete          = 'SnapshotDelete';
    case SnapshotExport          = 'SnapshotExport';
    case SnapshotImport          = 'SnapshotImport';
    case SnapshotCleanup         = 'SnapshotCleanup';
    case SnapshotFullBackup      = 'SnapshotFullBackup';
    case SnapshotIncremental     = 'SnapshotIncremental';
    case SnapshotSettingsUpdate  = 'SnapshotSettingsUpdate';
    case SnapshotZipBuild        = 'SnapshotZipBuild';
    case SnapshotZipExpire       = 'SnapshotZipExpire';
    case SnapshotZipDownload     = 'SnapshotZipDownload';

    // ── Cloud Storage Actions ───────────────────────────────
    case CloudStorageUpload        = 'CloudStorageUpload';
    case CloudStorageDelete        = 'CloudStorageDelete';
    case CloudStorageRotation      = 'CloudStorageRotation';
    case CloudStorageAccountAdd    = 'CloudStorageAccountAdd';
    case CloudStorageAccountRemove = 'CloudStorageAccountRemove';

    // ── Metadata ───────────────────────────────────────────

    public function label(): string
    {
        return match ($this) {
            self::Upload           => 'Plugin uploaded',
            self::UploadActive     => 'Active plugin uploaded',
            self::UploadInitiated  => 'Upload initiated',
            self::Enable           => 'Plugin enabled',
            self::Disable          => 'Plugin disabled',
            self::Delete           => 'Plugin deleted',
            self::FileReplace      => 'File replaced',
            self::FileDelete       => 'File deleted',
            self::Sync             => 'Sync executed',
            self::SyncDelete       => 'Sync delete executed',
            self::PostCreate       => 'Post created',
            self::PostUpdate       => 'Post updated',
            self::CategoryCreate   => 'Category created',
            self::MediaUpload      => 'Media uploaded',
            self::AuthFailed       => 'Authentication failed',
            self::ExportSelf       => 'Self-export executed',
            self::ExportPlugin     => 'Plugin export executed',
            self::PluginBackup        => 'Plugin backup created',
            self::PluginBackupRestore => 'Plugin backup restored',
            self::PluginBackupDelete  => 'Plugin backup deleted',
            self::AgentAdd           => 'Agent added',
            self::AgentRemove        => 'Agent removed',
            self::AgentTest          => 'Agent connection tested',
            self::AgentSync          => 'Agent sync executed',
            self::AgentApiError      => 'Agent API error',
            self::SnapshotCreate          => 'Snapshot created',
            self::SnapshotRestore         => 'Snapshot restored',
            self::SnapshotDelete          => 'Snapshot deleted',
            self::SnapshotExport          => 'Snapshot exported',
            self::SnapshotImport          => 'Snapshot imported',
            self::SnapshotCleanup         => 'Snapshot cleanup executed',
            self::SnapshotFullBackup      => 'Full snapshot backup created',
            self::SnapshotIncremental     => 'Incremental snapshot created',
            self::SnapshotSettingsUpdate  => 'Snapshot settings updated',
            self::SnapshotZipBuild        => 'Snapshot ZIP built',
            self::SnapshotZipExpire       => 'Snapshot ZIP expired',
            self::SnapshotZipDownload     => 'Snapshot ZIP downloaded',
            self::CloudStorageUpload        => 'Cloud storage upload',
            self::CloudStorageDelete        => 'Cloud storage file deleted',
            self::CloudStorageRotation      => 'Cloud storage rotation executed',
            self::CloudStorageAccountAdd    => 'Cloud storage account added',
            self::CloudStorageAccountRemove => 'Cloud storage account removed',
        };
    }

    // ── Per-Case Helpers ────────────────────────────────────

    public function isUpload(): bool           { return $this->isEqual(self::Upload); }
    public function isUploadActive(): bool     { return $this->isEqual(self::UploadActive); }
    public function isUploadInitiated(): bool  { return $this->isEqual(self::UploadInitiated); }
    public function isEnable(): bool           { return $this->isEqual(self::Enable); }
    public function isDisable(): bool          { return $this->isEqual(self::Disable); }
    public function isDelete(): bool           { return $this->isEqual(self::Delete); }
    public function isFileReplace(): bool      { return $this->isEqual(self::FileReplace); }
    public function isFileDelete(): bool       { return $this->isEqual(self::FileDelete); }
    public function isSync(): bool             { return $this->isEqual(self::Sync); }
    public function isSyncDelete(): bool       { return $this->isEqual(self::SyncDelete); }
    public function isPostCreate(): bool       { return $this->isEqual(self::PostCreate); }
    public function isPostUpdate(): bool       { return $this->isEqual(self::PostUpdate); }
    public function isCategoryCreate(): bool   { return $this->isEqual(self::CategoryCreate); }
    public function isMediaUpload(): bool      { return $this->isEqual(self::MediaUpload); }
    public function isAuthFailed(): bool       { return $this->isEqual(self::AuthFailed); }
    public function isExportSelf(): bool       { return $this->isEqual(self::ExportSelf); }
    public function isExportPlugin(): bool     { return $this->isEqual(self::ExportPlugin); }
    public function isPluginBackup(): bool        { return $this->isEqual(self::PluginBackup); }
    public function isPluginBackupRestore(): bool { return $this->isEqual(self::PluginBackupRestore); }
    public function isPluginBackupDelete(): bool  { return $this->isEqual(self::PluginBackupDelete); }
    public function isAgentAdd(): bool           { return $this->isEqual(self::AgentAdd); }
    public function isAgentRemove(): bool        { return $this->isEqual(self::AgentRemove); }
    public function isAgentTest(): bool          { return $this->isEqual(self::AgentTest); }
    public function isAgentSync(): bool          { return $this->isEqual(self::AgentSync); }
    public function isAgentApiError(): bool      { return $this->isEqual(self::AgentApiError); }
    public function isSnapshotCreate(): bool          { return $this->isEqual(self::SnapshotCreate); }
    public function isSnapshotRestore(): bool         { return $this->isEqual(self::SnapshotRestore); }
    public function isSnapshotDelete(): bool          { return $this->isEqual(self::SnapshotDelete); }
    public function isSnapshotExport(): bool          { return $this->isEqual(self::SnapshotExport); }
    public function isSnapshotImport(): bool          { return $this->isEqual(self::SnapshotImport); }
    public function isSnapshotCleanup(): bool         { return $this->isEqual(self::SnapshotCleanup); }
    public function isSnapshotFullBackup(): bool      { return $this->isEqual(self::SnapshotFullBackup); }
    public function isSnapshotIncremental(): bool     { return $this->isEqual(self::SnapshotIncremental); }
    public function isSnapshotSettingsUpdate(): bool  { return $this->isEqual(self::SnapshotSettingsUpdate); }
    public function isSnapshotZipBuild(): bool        { return $this->isEqual(self::SnapshotZipBuild); }
    public function isSnapshotZipExpire(): bool       { return $this->isEqual(self::SnapshotZipExpire); }
    public function isSnapshotZipDownload(): bool     { return $this->isEqual(self::SnapshotZipDownload); }
    public function isCloudStorageUpload(): bool        { return $this->isEqual(self::CloudStorageUpload); }
    public function isCloudStorageDelete(): bool        { return $this->isEqual(self::CloudStorageDelete); }
    public function isCloudStorageRotation(): bool      { return $this->isEqual(self::CloudStorageRotation); }
    public function isCloudStorageAccountAdd(): bool    { return $this->isEqual(self::CloudStorageAccountAdd); }
    public function isCloudStorageAccountRemove(): bool { return $this->isEqual(self::CloudStorageAccountRemove); }

    // ── Group Helpers ───────────────────────────────────────

    public function isSnapshot(): bool     { return str_starts_with($this->value, 'Snapshot'); }
    public function isAgent(): bool        { return str_starts_with($this->value, 'Agent'); }
    public function isCloudStorage(): bool { return str_starts_with($this->value, 'CloudStorage'); }

    public function isLifecycle(): bool
    {
        return $this->isAnyOf(self::Enable, self::Disable, self::Delete);
    }

    public function isExport(): bool
    {
        return $this->isAnyOf(self::ExportSelf, self::ExportPlugin);
    }

    public function isPluginBackupAction(): bool
    {
        return $this->isAnyOf(self::PluginBackup, self::PluginBackupRestore, self::PluginBackupDelete);
    }

    // ── Standard Comparison Methods ─────────────────────────

    public function isEqual(self $other): bool { return $this === $other; }
    public function isOtherThan(self $other): bool { return $this !== $other; }
    public function isAnyOf(self ...$others): bool { return in_array($this, $others, true); }
}
```

---

## Usage

```php
use PluginName\Enums\ActionType;

$action = ActionType::SnapshotZipExpire;

// Per-case helper
$action->isSnapshotZipExpire();  // true
$action->isUpload();             // false

// Group helper
$action->isSnapshot();           // true
$action->isLifecycle();          // false

// Label via match
$action->label();
// → "Snapshot ZIP expired"
```

---

## Cross-References

- [02-enum-metadata-pattern.md](02-enum-metadata-pattern.md) — the pattern specification
- [03-self-update-status-enum.md](03-self-update-status-enum.md) — first reference implementation (17 cases)
- [01-enum-architecture.md](01-enum-architecture.md) — core enum rules and comparison methods

---

*Second reference implementation — large enum (40+ cases) with match-based metadata.*
