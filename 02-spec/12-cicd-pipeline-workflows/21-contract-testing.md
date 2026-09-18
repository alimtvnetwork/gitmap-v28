# Contract Testing Standards

## 1. Purpose

Contract testing ensures that independent services (e.g., frontend and backend microservices) can communicate correctly without requiring fragile and slow end-to-end (E2E) integration environments.

## 2. Consumer-Driven Contracts

- We follow the Consumer-Driven Contract (CDC) pattern using frameworks like **Pact**.
- The Consumer (e.g., Frontend) defines the contract: the exact request it makes and the exact response it expects.
- The Provider (e.g., Backend) verifies that it can fulfill this contract.

## 3. Pipeline Integration

- Changes to Consumer expectations must publish a new contract to the Pact Broker.
- Provider CI pipelines must download the latest contracts and verify them before allowing a merge or deployment.
- Breaking a contract is treated as a CODE-RED failure in the Provider's pipeline.
