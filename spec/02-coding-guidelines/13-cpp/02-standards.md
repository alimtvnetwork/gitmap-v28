# Modern C++ Standards

## 1. C++20 Baseline

- All code MUST be written against the C++20 standard or newer.
- Use concepts instead of raw `enable_if` for template constraints.

## 2. Memory Safety & RAII

- **Never use `new` or `delete` manually.** Always use `std::unique_ptr` or `std::shared_ptr`.
- Follow the **Rule of Five** (or Rule of Zero preferred). If a class defines a custom destructor, copy constructor, move constructor, copy assignment, or move assignment, it should explicitly define or delete all five.
- Resources must be tied to object lifecycles (RAII).

## 3. Naming Conventions

- Structs, Classes, and Enums: `PascalCase`
- Functions and Variables: `snake_case` (follows community standards similar to Rust)
- Macros: `SCREAMING_SNAKE_CASE` (avoid macros where `constexpr` can be used)

## 4. Exceptions

- Avoid exceptions for control flow.
- When crossing C/C++ boundaries, ensure exceptions do not escape C++ code.
