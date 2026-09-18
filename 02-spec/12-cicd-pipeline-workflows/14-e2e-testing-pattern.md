# End-to-End (E2E) Testing Patterns

## 1. Page Object Model (POM)

- All E2E tests (Playwright, Cypress) MUST utilize the Page Object Model design pattern.
- Test files must never contain direct DOM selectors (e.g., `cy.get('.btn-primary')`).
- DOM selectors and interaction logic must be encapsulated within a Page Class (e.g., `LoginPage.ts`). Test files simply call methods on these classes (e.g., `LoginPage.login(user, pass)`).

## 2. Data Attributes for Selectors

- E2E tests should rely on dedicated `data-testid` or `data-cy` attributes rather than brittle CSS classes or element IDs.

## 3. Test Isolation

- Each E2E test must be completely isolated and capable of running in parallel.
- Tests must clean up their own state or rely on database transaction rollbacks. Never rely on the state left over by a previous test.
