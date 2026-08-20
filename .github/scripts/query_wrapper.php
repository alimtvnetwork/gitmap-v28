<?php

/**
 * Generic query wrapper result structure.
 */
class QueryResult {
    public $data;
    public $error;
    public bool $isSuccess;
    public bool $isFailure;

    public function __construct($data, $error, bool $isSuccess, bool $isFailure) {
        $this->data = $data;
        $this->error = $error;
        $this->isSuccess = $isSuccess;
        $this->isFailure = $isFailure;
    }
}

/**
 * Wraps an operation (API, database, or command logic), catching exceptions,
 * logging errors explicitly to stderr without scattered try/catch, 
 * and returning a structured result with explicit `isSuccess` and `isFailure` properties.
 */
function query_wrapper(callable $operation, ...$args): QueryResult {
    try {
        $data = call_user_func_array($operation, $args);
        return new QueryResult($data, null, true, false);
    } catch (Exception $e) {
        // Explicitly log the caught error to stderr
        file_put_contents('php://stderr', "[QueryWrapper Error]: " . $e->getMessage() . PHP_EOL);
        
        return new QueryResult(null, $e, false, true);
    }
}
