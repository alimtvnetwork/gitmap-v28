# Verbatim Directive: Version Architecture and Code Quality

- **Date**: 2026-08-09
- **Topic**: Version injection, Code quality enforcement

> "Fix the git status first, then start coding. Make a big plan if required to self-loop... Look into the entire codebase and follow the code review guidelines from the aspect folder properly. All caught errors must be explicitly logged following the guidelines in the error manage folder. Create a wrapper for queries in PHP/Python/TS that automatically logs failures to reduce scattered logging code.
> 
> Make sure the code quality is strictly maintained:
> 1. Do not introduce any magic strings or magic numbers anywhere unless it is explicitly for the logger, and mention that in the typing.
> 2. In TypeScript, rather than using strings as sub-items or comparing string union types (pipes) like "pass" | "fail" | "fallback", you must use Enums. Enums are the best.
> 3. Every single Enum must end with the suffix "Type".
> 4. Always use explicit boolean state checks like response.isFail or explicit checks rather than inverting success booleans like !response.isSuccess.
> ...
> Finally when your tasks are done, make sure you made a final bump in the minor release with following proper steps of release for this repo... Git should be source of the truth."

> "So when you, uh, write the code, you should not ask it to touch the gitmap folder. Gitmap folder is out of the touch... there should be a version of JSON file, which actually contains the version information. From this, every one of the version from, for let's say Golang or anything, should read that and then compile the code. I think that should be the easiest way to go. And based on that, you should actually update your MD files according to this... gitmap, gitmap, uh, latest. Um, inside the plan, you don't also need to add the version information. The version should be very simple to update only using that version.json file, which is root of the repo. And every one of the code base should follow that JSON file change... update inside the zero seven folder, zero eight, and zero nine as well."
