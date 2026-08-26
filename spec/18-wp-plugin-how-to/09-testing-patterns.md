# Phase 9 — Testing Patterns

> **Purpose:** Define how to write unit and integration tests for WordPress plugins built with this architecture. Covers testing traits, helpers, enums, the EnvelopeBuilder, FileLogger, and REST endpoints — all without requiring a running WordPress instance for unit tests.

---

## 9.1 Testing Philosophy

| Principle | Rule |
|-----------|------|
| **Unit tests don't need WordPress** | Helpers, enums, and pure logic are tested standalone with PHPUnit |
| **Integration tests use WP test suite** | REST endpoints, hooks, and database are tested with `wp-phpunit` |
| **Every public method has a test** | Enums, helpers, and trait public methods are all covered |
| **Test names describe behaviour** | `testValidationRejectsEmptyName` not `testValidation1` |
| **No test depends on another** | Each test sets up its own state and tears it down |

---

## 9.2 Test Directory Structure

```
plugin-slug/
├── tests/
│   ├── bootstrap.php              ← Test bootstrap (loads autoloader, mocks ABSPATH)
│   ├── Unit/
│   │   ├── Enums/
│   │   │   ├── PluginConfigTypeTest.php
│   │   │   ├── HttpStatusTypeTest.php
│   │   │   ├── PhpNativeTypeTest.php
│   │   │   └── ResponseKeyTypeTest.php
│   │   ├── Helpers/
│   │   │   ├── EnvelopeBuilderTest.php
│   │   │   ├── DateHelperTest.php
│   │   │   └── PathHelperTest.php
│   │   ├── Traits/
│   │   │   ├── TypeCheckerTraitTest.php
│   │   │   └── ResponseTraitTest.php
│   │   └── Database/
│   │       └── DatabaseSeederTest.php
│   ├── Integration/
│   │   ├── Endpoints/
│   │   │   ├── StatusEndpointTest.php
│   │   │   └── UploadEndpointTest.php
│   │   └── Cron/
│   │       └── CronSchedulerTest.php
│   └── Fixtures/
│       ├── sample-upload.zip
│       ├── invalid-file.txt
│       └── seeds/                   ← Test seed fixtures
│           ├── manifest.json
│           └── settings.json
├── phpunit.xml
└── composer.json                  ← PHPUnit + wp-phpunit as dev deps
```

---

## 9.3 Test Bootstrap — `tests/bootstrap.php`

For **unit tests**, the bootstrap mocks WordPress constants and functions so tests run without WordPress:

```php
<?php
/**
 * Test bootstrap — loads autoloader and mocks WordPress environment.
 */

// Define ABSPATH so guarded files load
if (!defined('ABSPATH')) {
    define('ABSPATH', '/tmp/fake-wordpress/');
}

// Mock WordPress functions used by helpers/enums (only for unit tests)
if (!function_exists('sanitize_text_field')) {
    function sanitize_text_field(string $str): string
    {
        return trim(strip_tags($str));
    }
}

if (!function_exists('absint')) {
    function absint($value): int
    {
        return abs((int) $value);
    }
}

if (!function_exists('wp_upload_dir')) {
    function wp_upload_dir(): array
    {
        return ['basedir' => sys_get_temp_dir() . '/wp-uploads'];
    }
}

if (!function_exists('wp_mkdir_p')) {
    function wp_mkdir_p(string $path): bool
    {
        $exists = is_dir($path);

        if ($exists) {
            return true;
        }

        return mkdir($path, 0755, true);
    }
}

// Load the plugin autoloader
require_once __DIR__ . '/../includes/Autoloader.php';
```

### Why mock WordPress functions?

Unit tests must run in CI without a WordPress installation. Only mock the functions your tested code actually calls. Keep mocks minimal — complex behaviour belongs in integration tests.

---

## 9.4 PHPUnit Configuration — `phpunit.xml`

```xml
<?xml version="1.0" encoding="UTF-8"?>
<phpunit
    bootstrap="tests/bootstrap.php"
    colors="true"
    stopOnFailure="false"
    cacheDirectory=".phpunit.cache"
>
    <testsuites>
        <testsuite name="Unit">
            <directory>tests/Unit</directory>
        </testsuite>
        <testsuite name="Integration">
            <directory>tests/Integration</directory>
        </testsuite>
    </testsuites>

    <source>
        <include>
            <directory>includes</directory>
        </include>
        <exclude>
            <file>includes/Autoloader.php</file>
        </exclude>
    </source>
</phpunit>
```

Run unit tests only: `vendor/bin/phpunit --testsuite Unit`
Run everything: `vendor/bin/phpunit`

---

## 9.5 Testing Enums

Enums are the easiest to test — pure value objects with no dependencies.

### PluginConfigTypeTest.php

```php
<?php

namespace PluginName\Tests\Unit\Enums;

use PHPUnit\Framework\TestCase;
use PluginName\Enums\PluginConfigType;

final class PluginConfigTypeTest extends TestCase
{
    public function testSlugIsKebabCase(): void
    {
        $slug = PluginConfigType::Slug->value;
        $isKebabCase = (preg_match('/^[a-z0-9]+(-[a-z0-9]+)*$/', $slug) === 1);

        $this->assertTrue($isKebabCase, "Slug '{$slug}' must be kebab-case");
    }

    public function testVersionFollowsSemver(): void
    {
        $version = PluginConfigType::Version->value;
        $isSemver = (preg_match('/^\d+\.\d+\.\d+$/', $version) === 1);

        $this->assertTrue($isSemver, "Version '{$version}' must follow semver (x.y.z)");
    }

    public function testApiFullNamespaceFormat(): void
    {
        $namespace = PluginConfigType::apiFullNamespace();

        $this->assertStringContainsString('/', $namespace);
        $this->assertStringEndsWith('/v1', $namespace);
    }

    public function testDebugModeReturnsBool(): void
    {
        $result = PluginConfigType::isDebugMode();

        $this->assertIsBool($result);
    }

    public function testIsEqualComparison(): void
    {
        $slug = PluginConfigType::Slug;

        $this->assertTrue($slug->isEqual(PluginConfigType::Slug));
        $this->assertFalse($slug->isEqual(PluginConfigType::Version));
    }

    public function testIsOtherThanComparison(): void
    {
        $slug = PluginConfigType::Slug;

        $this->assertTrue($slug->isOtherThan(PluginConfigType::Version));
        $this->assertFalse($slug->isOtherThan(PluginConfigType::Slug));
    }

    public function testIsAnyOfComparison(): void
    {
        $slug = PluginConfigType::Slug;

        $this->assertTrue($slug->isAnyOf(PluginConfigType::Slug, PluginConfigType::Name));
        $this->assertFalse($slug->isAnyOf(PluginConfigType::Version, PluginConfigType::Name));
    }
}
```

### PhpNativeTypeTest.php

```php
<?php

namespace PluginName\Tests\Unit\Enums;

use PHPUnit\Framework\TestCase;
use PluginName\Enums\PhpNativeType;

final class PhpNativeTypeTest extends TestCase
{
    /**
     * @dataProvider matchesProvider
     */
    public function testMatchesReturnsCorrectResult(
        PhpNativeType $type,
        mixed $value,
        bool $expected,
    ): void {
        $this->assertSame($expected, $type->matches($value));
    }

    /**
     * @return array<string, array{PhpNativeType, mixed, bool}>
     */
    public static function matchesProvider(): array
    {
        return [
            'array matches array'       => [PhpNativeType::PhpArray, [1, 2], true],
            'array rejects string'      => [PhpNativeType::PhpArray, 'hello', false],
            'string matches string'     => [PhpNativeType::PhpString, 'hello', true],
            'string rejects int'        => [PhpNativeType::PhpString, 42, false],
            'integer matches int'       => [PhpNativeType::PhpInteger, 42, true],
            'integer rejects float'     => [PhpNativeType::PhpInteger, 3.14, false],
            'double matches float'      => [PhpNativeType::PhpDouble, 3.14, true],
            'double rejects int'        => [PhpNativeType::PhpDouble, 42, false],
            'boolean matches bool'      => [PhpNativeType::PhpBoolean, true, true],
            'boolean rejects string'    => [PhpNativeType::PhpBoolean, 'true', false],
            'null matches null'         => [PhpNativeType::PhpNull, null, true],
            'null rejects empty string' => [PhpNativeType::PhpNull, '', false],
            'object matches object'     => [PhpNativeType::PhpObject, new \stdClass(), true],
            'object rejects array'      => [PhpNativeType::PhpObject, [], false],
        ];
    }

    public function testAllCasesHaveValidGettypeValue(): void
    {
        $validTypes = ['array', 'string', 'integer', 'double', 'boolean', 'object', 'NULL'];

        foreach (PhpNativeType::cases() as $case) {
            $this->assertContains(
                $case->value,
                $validTypes,
                "Enum case {$case->name} has invalid gettype value: {$case->value}",
            );
        }
    }
}
```

### Pattern: Every enum test should verify

- [ ] All cases have the expected backing values
- [ ] `isEqual()`, `isOtherThan()`, `isAnyOf()` work correctly
- [ ] Helper methods return expected types
- [ ] Format constraints are met (kebab-case slugs, semver versions, etc.)

---

## 9.6 Testing the TypeCheckerTrait

Traits cannot be instantiated directly. Create a concrete test class that uses the trait:

```php
<?php

namespace PluginName\Tests\Unit\Traits;

use PHPUnit\Framework\TestCase;
use PluginName\Traits\Core\TypeCheckerTrait;

final class TypeCheckerTraitTest extends TestCase
{
    private object $checker;

    protected function setUp(): void
    {
        // Anonymous class that uses the trait
        $this->checker = new class {
            use TypeCheckerTrait {
                isArray as public;
                isString as public;
                isInteger as public;
                isFloat as public;
                isBoolean as public;
                isObject as public;
                isNull as public;
                isNumeric as public;
                isScalar as public;
            }
        };
    }

    public function testIsArrayWithArray(): void
    {
        $this->assertTrue($this->checker->isArray([]));
        $this->assertTrue($this->checker->isArray([1, 2, 3]));
        $this->assertTrue($this->checker->isArray(['key' => 'value']));
    }

    public function testIsArrayRejectsNonArrays(): void
    {
        $this->assertFalse($this->checker->isArray('not array'));
        $this->assertFalse($this->checker->isArray(42));
        $this->assertFalse($this->checker->isArray(null));
        $this->assertFalse($this->checker->isArray(new \stdClass()));
    }

    public function testIsStringWithString(): void
    {
        $this->assertTrue($this->checker->isString('hello'));
        $this->assertTrue($this->checker->isString(''));
    }

    public function testIsStringRejectsNonStrings(): void
    {
        $this->assertFalse($this->checker->isString(42));
        $this->assertFalse($this->checker->isString(true));
        $this->assertFalse($this->checker->isString(null));
    }

    public function testIsIntegerWithInt(): void
    {
        $this->assertTrue($this->checker->isInteger(0));
        $this->assertTrue($this->checker->isInteger(-1));
        $this->assertTrue($this->checker->isInteger(PHP_INT_MAX));
    }

    public function testIsIntegerRejectsFloatAndNumericString(): void
    {
        $this->assertFalse($this->checker->isInteger(3.14));
        $this->assertFalse($this->checker->isInteger('42'));
    }

    public function testIsNumericAcceptsBothIntAndFloat(): void
    {
        $this->assertTrue($this->checker->isNumeric(42));
        $this->assertTrue($this->checker->isNumeric(3.14));
        $this->assertFalse($this->checker->isNumeric('42'));
        $this->assertFalse($this->checker->isNumeric(null));
    }

    public function testIsScalarAcceptsAllScalarTypes(): void
    {
        $this->assertTrue($this->checker->isScalar('hello'));
        $this->assertTrue($this->checker->isScalar(42));
        $this->assertTrue($this->checker->isScalar(3.14));
        $this->assertTrue($this->checker->isScalar(true));
    }

    public function testIsScalarRejectsCompoundTypes(): void
    {
        $this->assertFalse($this->checker->isScalar([]));
        $this->assertFalse($this->checker->isScalar(new \stdClass()));
        $this->assertFalse($this->checker->isScalar(null));
    }

    public function testIsNullOnlyMatchesNull(): void
    {
        $this->assertTrue($this->checker->isNull(null));
        $this->assertFalse($this->checker->isNull(false));
        $this->assertFalse($this->checker->isNull(0));
        $this->assertFalse($this->checker->isNull(''));
    }
}
```

### Key technique: Trait visibility override

```php
new class {
    use TypeCheckerTrait {
        isArray as public;  // Override protected → public for testing
    }
};
```

This is the standard PHPUnit pattern for testing `protected` trait methods without creating a full concrete class.

---

## 9.7 Testing the EnvelopeBuilder

```php
<?php

namespace PluginName\Tests\Unit\Helpers;

use PHPUnit\Framework\TestCase;
use PluginName\Helpers\EnvelopeBuilder;

final class EnvelopeBuilderTest extends TestCase
{
    // ── Success Responses ──

    public function testSuccessEnvelopeHasCorrectStructure(): void
    {
        $response = EnvelopeBuilder::success('OK', 200)
            ->setRequestedAt('/test/v1/endpoint')
            ->setSingleResult(['id' => '123'])
            ->toResponse();

        $data = $response->get_data();

        $this->assertSame(200, $response->get_status());
        $this->assertTrue($data['Status']['IsSuccess']);
        $this->assertFalse($data['Status']['IsFailed']);
        $this->assertSame('OK', $data['Status']['Message']);
        $this->assertSame('/test/v1/endpoint', $data['Attributes']['RequestedAt']);
        $this->assertCount(1, $data['Results']);
        $this->assertSame('123', $data['Results'][0]['id']);
        $this->assertArrayNotHasKey('Errors', $data);
    }

    public function testSuccessWithListResult(): void
    {
        $items = [['id' => '1'], ['id' => '2'], ['id' => '3']];

        $response = EnvelopeBuilder::success()
            ->setListResult($items)
            ->toResponse();

        $data = $response->get_data();

        $this->assertSame(3, $data['Attributes']['TotalRecords']);
        $this->assertCount(3, $data['Results']);
    }

    public function testSuccessWithEmptyResults(): void
    {
        $response = EnvelopeBuilder::success()->toResponse();
        $data = $response->get_data();

        $this->assertSame(0, $data['Attributes']['TotalRecords']);
        $this->assertSame([], $data['Results']);
    }

    // ── Error Responses ──

    public function testErrorEnvelopeWithoutDebugOmitsStackTrace(): void
    {
        // PLUGIN_NAME_DEBUG is false (or undefined) in test environment
        $exception = new \RuntimeException('Something broke');

        $response = EnvelopeBuilder::error('Something broke', 500, $exception)
            ->setRequestedAt('/test/v1/endpoint')
            ->toResponse();

        $data = $response->get_data();

        $this->assertSame(500, $response->get_status());
        $this->assertFalse($data['Status']['IsSuccess']);
        $this->assertTrue($data['Status']['IsFailed']);

        // In production mode: generic message, no Errors key
        $this->assertSame('An internal error occurred', $data['Status']['Message']);
        $this->assertArrayNotHasKey('Errors', $data);
    }

    public function testErrorEnvelopeFor400IncludesRealMessage(): void
    {
        $response = EnvelopeBuilder::error('Missing field: name', 400)
            ->toResponse();

        $data = $response->get_data();

        // 400 errors always show real message regardless of debug mode
        $this->assertSame('Missing field: name', $data['Status']['Message']);
    }

    // ── Timestamp ──

    public function testTimestampIsUtcIso8601(): void
    {
        $response = EnvelopeBuilder::success()->toResponse();
        $data = $response->get_data();

        $timestamp = $data['Status']['Timestamp'];
        $parsed = \DateTimeImmutable::createFromFormat(\DateTimeInterface::ATOM, $timestamp);
        $isValidTimestamp = ($parsed !== false);

        $this->assertTrue($isValidTimestamp, "Timestamp '{$timestamp}' is not valid ISO 8601");
    }

    // ── PascalCase Keys ──

    public function testAllTopLevelKeysArePascalCase(): void
    {
        $response = EnvelopeBuilder::success()
            ->setSingleResult(['test' => true])
            ->toResponse();

        $data = $response->get_data();
        $keys = array_keys($data);

        foreach ($keys as $key) {
            $isPascalCase = (preg_match('/^[A-Z][a-zA-Z]*$/', $key) === 1);

            $this->assertTrue($isPascalCase, "Key '{$key}' is not PascalCase");
        }
    }
}
```

### Testing debug-mode responses

To test the debug-mode path, define the constant before the test:

```php
public function testErrorEnvelopeWithDebugIncludesStackTrace(): void
{
    // This test only works if PLUGIN_NAME_DEBUG is true
    // Define it in a separate phpunit.xml bootstrap or test-specific setup
    if (!defined('PLUGIN_NAME_DEBUG')) {
        define('PLUGIN_NAME_DEBUG', true);
    }

    $exception = new \RuntimeException('Test exception');

    $response = EnvelopeBuilder::error('Test exception', 500, $exception)
        ->toResponse();

    $data = $response->get_data();

    $this->assertArrayHasKey('Errors', $data);
    $this->assertSame('Test exception', $data['Errors']['BackendMessage']);
    $this->assertSame('RuntimeException', $data['Errors']['ExceptionType']);
    $this->assertNotEmpty($data['Errors']['Backend']);
}
```

> **Warning:** PHP constants can only be defined once per process. If you need to test both debug ON and OFF, use separate PHPUnit test suites with different bootstrap files, or use a method-based approach in `PluginConfigType::isDebugMode()` that can be overridden in tests.

### Testable debug mode pattern

To make debug mode testable without constants:

```php
// In PluginConfigType — add override for tests
private static ?bool $debugOverride = null;

public static function setDebugOverride(?bool $value): void
{
    self::$debugOverride = $value;
}

public static function isDebugMode(): bool
{
    $hasOverride = (self::$debugOverride !== null);

    if ($hasOverride) {
        return self::$debugOverride;
    }

    $constantName = self::DebugConstant->value;
    $isDefined = defined($constantName);

    return $isDefined && constant($constantName) === true;
}
```

Then in tests:

```php
protected function setUp(): void
{
    PluginConfigType::setDebugOverride(true);
}

protected function tearDown(): void
{
    PluginConfigType::setDebugOverride(null);
}
```

---

## 9.8 Testing REST Endpoints (Integration)

Integration tests require the WordPress test suite (`wp-phpunit`). They test the full request-response cycle.

### composer.json dev dependencies

```json
{
    "require-dev": {
        "phpunit/phpunit": "^10.0",
        "yoast/wp-test-utils": "^1.0"
    }
}
```

### Integration test bootstrap

```php
<?php
/**
 * Integration test bootstrap — loads WordPress test suite.
 */

$testsDir = getenv('WP_TESTS_DIR') ?: '/tmp/wordpress-tests-lib';

require_once $testsDir . '/includes/functions.php';

// Load the plugin
tests_add_filter('muplugins_loaded', function (): void {
    require dirname(__DIR__, 2) . '/plugin-name.php';
});

require $testsDir . '/includes/bootstrap.php';
```

### StatusEndpointTest.php

```php
<?php

namespace PluginName\Tests\Integration\Endpoints;

use WP_REST_Request;
use WP_REST_Response;
use WP_UnitTestCase;

final class StatusEndpointTest extends WP_UnitTestCase
{
    private int $adminUserId;

    public function setUp(): void
    {
        parent::setUp();

        // Create an admin user for authenticated requests
        $this->adminUserId = $this->factory->user->create([
            'role' => 'administrator',
        ]);
    }

    public function testStatusEndpointReturnsSuccessEnvelope(): void
    {
        wp_set_current_user($this->adminUserId);

        $request = new WP_REST_Request('GET', '/plugin-name-api/v1/status');
        $response = rest_do_request($request);

        $this->assertSame(200, $response->get_status());

        $data = $response->get_data();
        $this->assertTrue($data['Status']['IsSuccess']);
        $this->assertArrayHasKey('Results', $data);
    }

    public function testStatusEndpointRejectsUnauthenticated(): void
    {
        wp_set_current_user(0); // No user

        $request = new WP_REST_Request('GET', '/plugin-name-api/v1/status');
        $response = rest_do_request($request);

        $status = $response->get_status();
        $isUnauthorized = ($status === 401 || $status === 403);

        $this->assertTrue($isUnauthorized);
    }

    public function testStatusEndpointRejectsNonAdmin(): void
    {
        $subscriberId = $this->factory->user->create(['role' => 'subscriber']);
        wp_set_current_user($subscriberId);

        $request = new WP_REST_Request('GET', '/plugin-name-api/v1/status');
        $response = rest_do_request($request);

        $this->assertSame(403, $response->get_status());
    }
}
```

---

## 9.9 Testing Validation (Unit)

Test the validation patterns from Phase 6 by creating a mock handler:

```php
<?php

namespace PluginName\Tests\Unit\Traits;

use PHPUnit\Framework\TestCase;

final class ValidationPatternTest extends TestCase
{
    private object $handler;

    protected function setUp(): void
    {
        $this->handler = new class {
            use \PluginName\Traits\Core\TypeCheckerTrait {
                isArray as public;
                isString as public;
                isInteger as public;
            }

            /**
             * Simulate the validation pattern for a "create widget" request.
             *
             * @param array<string, mixed>|null $body
             *
             * @return array{valid: bool, error: string|null, data: array<string, mixed>}
             */
            public function validateCreateWidget(?array $body): array
            {
                $hasBody = ($body !== null && $this->isArray($body));

                if (!$hasBody) {
                    return ['valid' => false, 'error' => 'Request body must be a JSON object', 'data' => []];
                }

                $name = $body['name'] ?? null;
                $hasName = ($name !== null && $this->isString($name));

                if (!$hasName) {
                    return ['valid' => false, 'error' => 'Missing required field: name', 'data' => []];
                }

                $nameLength = mb_strlen($name);
                $isNameTooLong = ($nameLength > 200);

                if ($isNameTooLong) {
                    return ['valid' => false, 'error' => 'Field "name" must not exceed 200 characters', 'data' => []];
                }

                return ['valid' => true, 'error' => null, 'data' => ['name' => $name]];
            }
        };
    }

    public function testRejectsNullBody(): void
    {
        $result = $this->handler->validateCreateWidget(null);

        $this->assertFalse($result['valid']);
        $this->assertSame('Request body must be a JSON object', $result['error']);
    }

    public function testRejectsMissingName(): void
    {
        $result = $this->handler->validateCreateWidget(['other' => 'value']);

        $this->assertFalse($result['valid']);
        $this->assertStringContainsString('name', $result['error']);
    }

    public function testRejectsNonStringName(): void
    {
        $result = $this->handler->validateCreateWidget(['name' => 42]);

        $this->assertFalse($result['valid']);
    }

    public function testRejectsOverlongName(): void
    {
        $longName = str_repeat('a', 201);
        $result = $this->handler->validateCreateWidget(['name' => $longName]);

        $this->assertFalse($result['valid']);
        $this->assertStringContainsString('200', $result['error']);
    }

    public function testAcceptsValidInput(): void
    {
        $result = $this->handler->validateCreateWidget(['name' => 'My Widget']);

        $this->assertTrue($result['valid']);
        $this->assertNull($result['error']);
        $this->assertSame('My Widget', $result['data']['name']);
    }

    public function testAcceptsMaxLengthName(): void
    {
        $maxName = str_repeat('a', 200);
        $result = $this->handler->validateCreateWidget(['name' => $maxName]);

        $this->assertTrue($result['valid']);
    }
}
```

---

## 9.10 Testing Data Providers — Edge Case Coverage

Use PHPUnit data providers to test boundary conditions systematically:

```php
/**
 * @dataProvider edgeCaseInputProvider
 */
public function testIsArrayEdgeCases(mixed $input, bool $expected): void
{
    $this->assertSame($expected, $this->checker->isArray($input));
}

/**
 * @return array<string, array{mixed, bool}>
 */
public static function edgeCaseInputProvider(): array
{
    return [
        'empty array'           => [[], true],
        'indexed array'         => [[1, 2, 3], true],
        'associative array'     => [['key' => 'val'], true],
        'nested array'          => [[[]], true],
        'string "array"'        => ['array', false],
        'integer zero'          => [0, false],
        'boolean false'         => [false, false],
        'null'                  => [null, false],
        'empty string'          => ['', false],
        'stdClass object'       => [new \stdClass(), false],
        'ArrayObject'           => [new \ArrayObject(), false],  // Important: not a native array!
    ];
}
```

### Critical edge case: `ArrayObject` is not an array

`gettype(new ArrayObject()) === 'object'` — this catches a common mistake where code assumes array-like objects are arrays. The TypeCheckerTrait correctly rejects them.

---

## 9.11 Test Naming Convention

| Pattern | Example |
|---------|---------|
| `test{Action}{Condition}` | `testValidationRejectsEmptyName` |
| `test{Subject}{Behaviour}` | `testEnvelopeBuilderOmitsErrorsOnSuccess` |
| `test{Subject}{EdgeCase}` | `testIsArrayRejectsArrayObject` |

### Rules

- Test method names are `camelCase` starting with `test`
- Describe **what** is tested and **what** the expected outcome is
- Never use numeric suffixes (`testCase1`, `testCase2`)
- Group related tests with `// ── Section ──` comments

---

## 9.12 Coverage Targets

| Category | Minimum coverage | Notes |
|----------|-----------------|-------|
| Enums | 100% | All cases, all helper methods |
| TypeCheckerTrait | 100% | All type methods, all edge cases |
| EnvelopeBuilder | 95%+ | Both success and error paths, debug ON/OFF |
| Validation patterns | 90%+ | All field types, boundary values |
| Helpers (DateHelper, PathHelper) | 90%+ | All public methods |
| REST endpoints (integration) | 80%+ | Happy path + auth rejection + validation failure |
| FileLogger | 70%+ | Log levels, rotation trigger, dedup |

### Running with coverage

```bash
vendor/bin/phpunit --testsuite Unit --coverage-text --coverage-html coverage/
```

---

## 9.13 CI Integration

### GitHub Actions workflow

```yaml
name: PHPUnit
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        php: ['8.1', '8.2', '8.3']

    steps:
      - uses: actions/checkout@v4

      - name: Setup PHP
        uses: shivammathur/setup-php@v2
        with:
          php-version: ${{ matrix.php }}
          extensions: sqlite3
          coverage: xdebug

      - name: Install dependencies
        run: composer install --no-interaction

      - name: Run unit tests
        run: vendor/bin/phpunit --testsuite Unit

      - name: Run integration tests
        if: matrix.php == '8.2'  # Only run integration once
        run: vendor/bin/phpunit --testsuite Integration
        env:
          WP_TESTS_DIR: /tmp/wordpress-tests-lib
```

---

## 9.14 Testing Database Seeding

> **Ref:** Phase 8 §8.5.1 — Database Seeding patterns

### 9.14.1 Test Fixtures

Create minimal seed fixtures under `tests/Fixtures/seeds/`:

**`tests/Fixtures/seeds/manifest.json`**
```json
{
    "version": "1.0.0",
    "seeds": [
        {
            "file": "settings.json",
            "table": "test_settings",
            "strategy": "insert_if_empty",
            "key_column": "setting_key"
        }
    ]
}
```

**`tests/Fixtures/seeds/settings.json`**
```json
[
    { "setting_key": "color_primary", "setting_value": "#3B82F6" },
    { "setting_key": "max_uploads",   "setting_value": "10" }
]
```

### 9.14.2 DatabaseSeederTest

```php
<?php

declare(strict_types=1);

namespace PluginName\Tests\Unit\Database;

use PHPUnit\Framework\TestCase;
use PluginName\Database\DatabaseSeeder;

final class DatabaseSeederTest extends TestCase
{
    private \SQLite3 $db;
    private string   $seedDir;

    protected function setUp(): void
    {
        parent::setUp();

        $this->db = new \SQLite3(':memory:');
        $this->db->exec('PRAGMA foreign_keys = ON');

        // Create the seed_history tracking table
        $this->db->exec('
            CREATE TABLE seed_history (
                id               INTEGER PRIMARY KEY AUTOINCREMENT,
                seed_file        TEXT    NOT NULL UNIQUE,
                last_seeded_ver  TEXT    NOT NULL,
                seeded_at        TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP
            )
        ');

        // Create a target table for seed data
        $this->db->exec('
            CREATE TABLE test_settings (
                id            INTEGER PRIMARY KEY AUTOINCREMENT,
                setting_key   TEXT NOT NULL UNIQUE,
                setting_value TEXT NOT NULL
            )
        ');

        $this->seedDir = __DIR__ . '/../../Fixtures/seeds';
    }

    protected function tearDown(): void
    {
        $this->db->close();
        parent::tearDown();
    }

    // ── insert_if_empty ─────────────────────────────────────

    public function testInsertIfEmptyPopulatesEmptyTable(): void
    {
        $seeder = new DatabaseSeeder($this->db, $this->seedDir);
        $seeder->seedAll('1.0.0');

        $result = $this->db->querySingle(
            "SELECT COUNT(*) FROM test_settings"
        );
        $this->assertSame(2, $result);
    }

    public function testInsertIfEmptySkipsNonEmptyTable(): void
    {
        // Pre-populate — seeder should leave it alone
        $this->db->exec(
            "INSERT INTO test_settings (setting_key, setting_value)
             VALUES ('existing_key', 'existing_value')"
        );

        $seeder = new DatabaseSeeder($this->db, $this->seedDir);
        $seeder->seedAll('1.0.0');

        $result = $this->db->querySingle(
            "SELECT COUNT(*) FROM test_settings"
        );
        $this->assertSame(1, $result, 'insert_if_empty must skip non-empty tables');
    }

    // ── upsert_by_key ───────────────────────────────────────

    public function testUpsertByKeyUpdatesExistingRows(): void
    {
        // Insert a row the seed will match on key_column
        $this->db->exec(
            "INSERT INTO test_settings (setting_key, setting_value)
             VALUES ('color_primary', '#000000')"
        );

        $seeder = $this->createSeederWithStrategy('upsert_by_key');
        $seeder->seedAll('1.0.0');

        $value = $this->db->querySingle(
            "SELECT setting_value FROM test_settings
             WHERE setting_key = 'color_primary'"
        );
        $this->assertSame('#3B82F6', $value, 'upsert must update existing row');
    }

    public function testUpsertByKeyInsertsNewRows(): void
    {
        // Table exists but has no matching key
        $this->db->exec(
            "INSERT INTO test_settings (setting_key, setting_value)
             VALUES ('unrelated_key', 'value')"
        );

        $seeder = $this->createSeederWithStrategy('upsert_by_key');
        $seeder->seedAll('1.0.0');

        $result = $this->db->querySingle(
            "SELECT COUNT(*) FROM test_settings"
        );
        // 1 existing + 2 from seed
        $this->assertSame(3, $result, 'upsert must insert non-existent rows');
    }

    // ── replace_all ─────────────────────────────────────────

    public function testReplaceAllTruncatesAndReinserts(): void
    {
        // Pre-populate with data that should be wiped
        $this->db->exec(
            "INSERT INTO test_settings (setting_key, setting_value)
             VALUES ('old_key', 'old_value')"
        );

        $seeder = $this->createSeederWithStrategy('replace_all');
        $seeder->seedAll('1.0.0');

        $result = $this->db->querySingle(
            "SELECT COUNT(*) FROM test_settings"
        );
        $this->assertSame(2, $result, 'replace_all must truncate then insert seed data');

        // Old row must be gone
        $old = $this->db->querySingle(
            "SELECT setting_value FROM test_settings
             WHERE setting_key = 'old_key'"
        );
        $this->assertNull($old, 'replace_all must remove pre-existing rows');
    }

    // ── Version tracking ────────────────────────────────────

    public function testSeedingRecordsVersionInHistory(): void
    {
        $seeder = new DatabaseSeeder($this->db, $this->seedDir);
        $seeder->seedAll('1.0.0');

        $version = $this->db->querySingle(
            "SELECT last_seeded_ver FROM seed_history
             WHERE seed_file = 'settings.json'"
        );
        $this->assertSame('1.0.0', $version);
    }

    public function testSameVersionDoesNotReseed(): void
    {
        $seeder = new DatabaseSeeder($this->db, $this->seedDir);
        $seeder->seedAll('1.0.0');

        // Modify data manually
        $this->db->exec(
            "UPDATE test_settings SET setting_value = 'CHANGED'
             WHERE setting_key = 'color_primary'"
        );

        // Run seeder again with same version
        $seeder->seedAll('1.0.0');

        $value = $this->db->querySingle(
            "SELECT setting_value FROM test_settings
             WHERE setting_key = 'color_primary'"
        );
        $this->assertSame(
            'CHANGED',
            $value,
            'Same version must NOT re-seed — user changes preserved'
        );
    }

    public function testNewerVersionTriggersReseed(): void
    {
        $seeder = new DatabaseSeeder($this->db, $this->seedDir);
        $seeder->seedAll('1.0.0');

        // Wipe so insert_if_empty can re-populate
        $this->db->exec('DELETE FROM test_settings');

        // Bump version
        $seeder->seedAll('1.1.0');

        $result = $this->db->querySingle(
            "SELECT COUNT(*) FROM test_settings"
        );
        $this->assertSame(2, $result, 'Newer version must trigger re-seed');

        $version = $this->db->querySingle(
            "SELECT last_seeded_ver FROM seed_history
             WHERE seed_file = 'settings.json'"
        );
        $this->assertSame('1.1.0', $version);
    }

    public function testOlderVersionDoesNotReseed(): void
    {
        $seeder = new DatabaseSeeder($this->db, $this->seedDir);
        $seeder->seedAll('2.0.0');

        $this->db->exec('DELETE FROM test_settings');

        // Attempt with older version
        $seeder->seedAll('1.0.0');

        $result = $this->db->querySingle(
            "SELECT COUNT(*) FROM test_settings"
        );
        $this->assertSame(
            0,
            $result,
            'Older version must NOT re-seed'
        );
    }

    // ── Edge cases ──────────────────────────────────────────

    public function testMissingSeedFileThrowsException(): void
    {
        $seeder = new DatabaseSeeder($this->db, '/nonexistent/path');

        $this->expectException(\RuntimeException::class);
        $seeder->seedAll('1.0.0');
    }

    public function testMalformedJsonThrowsException(): void
    {
        $tmpDir = sys_get_temp_dir() . '/bad_seeds_' . uniqid();
        mkdir($tmpDir, 0755, true);

        file_put_contents($tmpDir . '/manifest.json', json_encode([
            'version' => '1.0.0',
            'seeds'   => [[
                'file'       => 'bad.json',
                'table'      => 'test_settings',
                'strategy'   => 'insert_if_empty',
                'key_column' => 'setting_key',
            ]],
        ]));
        file_put_contents($tmpDir . '/bad.json', '{ invalid json !!!');

        $seeder = new DatabaseSeeder($this->db, $tmpDir);

        $this->expectException(\RuntimeException::class);
        $seeder->seedAll('1.0.0');

        // Cleanup
        array_map('unlink', glob($tmpDir . '/*'));
        rmdir($tmpDir);
    }

    public function testEmptySeedArrayInsertsNothing(): void
    {
        $tmpDir = sys_get_temp_dir() . '/empty_seeds_' . uniqid();
        mkdir($tmpDir, 0755, true);

        file_put_contents($tmpDir . '/manifest.json', json_encode([
            'version' => '1.0.0',
            'seeds'   => [[
                'file'       => 'empty.json',
                'table'      => 'test_settings',
                'strategy'   => 'insert_if_empty',
                'key_column' => 'setting_key',
            ]],
        ]));
        file_put_contents($tmpDir . '/empty.json', '[]');

        $seeder = new DatabaseSeeder($this->db, $tmpDir);
        $seeder->seedAll('1.0.0');

        $result = $this->db->querySingle(
            "SELECT COUNT(*) FROM test_settings"
        );
        $this->assertSame(0, $result, 'Empty seed array must insert nothing');

        array_map('unlink', glob($tmpDir . '/*'));
        rmdir($tmpDir);
    }

    // ── Helper ──────────────────────────────────────────────

    /**
     * Creates a seeder with a modified manifest that uses the given strategy.
     */
    private function createSeederWithStrategy(string $strategy): DatabaseSeeder
    {
        $tmpDir = sys_get_temp_dir() . '/seeds_' . $strategy . '_' . uniqid();
        mkdir($tmpDir, 0755, true);

        // Copy seed data file
        copy($this->seedDir . '/settings.json', $tmpDir . '/settings.json');

        // Write manifest with overridden strategy
        file_put_contents($tmpDir . '/manifest.json', json_encode([
            'version' => '1.0.0',
            'seeds'   => [[
                'file'       => 'settings.json',
                'table'      => 'test_settings',
                'strategy'   => $strategy,
                'key_column' => 'setting_key',
            ]],
        ]));

        return new DatabaseSeeder($this->db, $tmpDir);
    }
}
```

### 9.14.3 What Each Test Proves

| Test | Strategy | Assertion |
|------|----------|-----------|
| `testInsertIfEmptyPopulatesEmptyTable` | `insert_if_empty` | Seeds populate empty tables |
| `testInsertIfEmptySkipsNonEmptyTable` | `insert_if_empty` | Existing data is never overwritten |
| `testUpsertByKeyUpdatesExistingRows` | `upsert_by_key` | Matching `key_column` rows get updated values |
| `testUpsertByKeyInsertsNewRows` | `upsert_by_key` | Non-matching keys are inserted alongside existing data |
| `testReplaceAllTruncatesAndReinserts` | `replace_all` | Table is wiped and re-populated from seed JSON |
| `testSeedingRecordsVersionInHistory` | — | `seed_history` table tracks `last_seeded_ver` |
| `testSameVersionDoesNotReseed` | — | Identical version skips seeding (user changes preserved) |
| `testNewerVersionTriggersReseed` | — | Bumped version triggers re-seed |
| `testOlderVersionDoesNotReseed` | — | Downgraded version is ignored |
| `testMissingSeedFileThrowsException` | — | Missing `data/seeds/` directory throws `RuntimeException` |
| `testMalformedJsonThrowsException` | — | Invalid JSON in seed file throws `RuntimeException` |
| `testEmptySeedArrayInsertsNothing` | — | Empty `[]` seed file inserts zero rows |

---

## 9.15 Test Checklist for New Features

When adding a new feature endpoint, also add:

- [ ] Enum test: new cases have valid backing values
- [ ] Validation test: each required field rejected when missing
- [ ] Validation test: each field rejected when wrong type
- [ ] Validation test: boundary values (max length, min/max range)
- [ ] Success test: valid input returns correct envelope
- [ ] Auth test: unauthenticated request rejected
- [ ] Auth test: insufficient capability rejected
- [ ] Edge case test: empty body, null values, oversized input
- [ ] Integration test: full round-trip via `rest_do_request()`
- [ ] Seeding test: new seed files have `insert_if_empty` / `upsert_by_key` coverage
