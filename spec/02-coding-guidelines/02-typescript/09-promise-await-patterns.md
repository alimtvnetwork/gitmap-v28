# Promise & Await Patterns

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**Applies to:** TypeScript / JavaScript  
**Source:** Consolidated from `01-pre-code-review-guides/02-fe-guidelines.md`

---

## 1. Rule

**Do not return a Promise and then immediately await it.** If you `await` a function, the Promise is already resolved — returning `Promise` is redundant.

---

## 2. Anti-Pattern: Return Promise + Await

```typescript
// ❌ WRONG — returning a Promise then awaiting it is contradictory
async function getData(): Promise<Data> {
    return new Promise(async (resolve) => {
        const result = await fetchApi('/data');

        resolve(result);
    });
}

// We return a Promise, which means the caller doesn't wait.
// But inside we await, which means we DO wait.
// This is a logical contradiction.
```

```typescript
// ✅ CORRECT — just use async/await directly
async function getData(): Promise<Data> {
    const result = await fetchApi('/data');

    return result;
}
```

---

## 3. Independent Promises MUST Run in Parallel — 🔴 CODE RED

> **Severity: CODE RED** — Sequential `await` on independent promises is an **automatic rejection** in code review. No exceptions.

When multiple async operations are **not dependent on each other**, they **MUST** be executed in parallel using `Promise.all`. Sequential `await` on independent calls wastes time proportional to the number of calls — a 3-call sequence takes 3× longer than parallel execution.

### How to Decide

| Question | Sequential `await` | `Promise.all` |
|----------|-------------------|---------------|
| Does call B need the result of call A? | ✅ Use sequential | — |
| Are calls A, B, C all independent? | — | ✅ **MANDATORY** |
| Are some calls independent, some dependent? | — | ✅ Parallel the independent ones, then sequential for dependent |

### Examples

```typescript
// ❌ CODE RED — sequential awaits on independent calls
async function loadDashboard(): Promise<Dashboard> {
    const users = await fetchUsers();       // waits 200ms
    const posts = await fetchPosts();       // waits 200ms (after users done)
    const stats = await fetchStats();       // waits 200ms (after posts done)
    // Total: ~600ms ❌

    return { users, posts, stats };
}

// ✅ REQUIRED — parallel with Promise.all
async function loadDashboard(): Promise<Dashboard> {
    const [users, posts, stats] = await Promise.all([
        fetchUsers(),
        fetchPosts(),
        fetchStats(),
    ]);
    // Total: ~200ms ✅ (all run simultaneously)

    return { users, posts, stats };
}
```

### Mixed Dependencies (Some Parallel, Some Sequential)

```typescript
// ✅ CORRECT — parallel where possible, sequential only when dependent
async function loadUserDashboard(userId: string): Promise<UserDashboard> {
    // Step 1: Fetch user first (others depend on it)
    const user = await fetchUser(userId);

    // Step 2: These are independent of each other — run in parallel
    const [posts, settings, notifications] = await Promise.all([
        fetchUserPosts(user.id),
        fetchUserSettings(user.id),
        fetchNotifications(user.id),
    ]);

    return { user, posts, settings, notifications };
}
```

### React Query / TanStack Query

In React components, use multiple `useQuery` hooks — they automatically run in parallel:

```typescript
// ✅ CORRECT — React Query runs these in parallel automatically
function Dashboard() {
    const users = useQuery({ queryKey: ['users'], queryFn: fetchUsers });
    const posts = useQuery({ queryKey: ['posts'], queryFn: fetchPosts });
    const stats = useQuery({ queryKey: ['stats'], queryFn: fetchStats });

    if (users.isLoading || posts.isLoading || stats.isLoading) {
        return <Skeleton />;
    }

    return <DashboardView users={users.data} posts={posts.data} stats={stats.data} />;
}
```

---

## 4. Error Handling with Async

```typescript
// ✅ CORRECT — structured error handling
async function fetchData(): Promise<Result<Data>> {
    try {
        const response = await apiClient.get('/data');

        if (!response.ok) {
            return { error: new AppError('Fetch failed', response.status) };
        }

        return { data: response.data };
    } catch (error) {
        return { error: new AppError('Network error', 500) };
    }
}
```

---

## 5. References

- [Promises or Async-Await (Better Programming)](https://betterprogramming.pub/should-i-use-promises-or-async-await-126ab5c98789)
- [Difference Between Promise and Async Await (GeeksForGeeks)](https://www.geeksforgeeks.org/difference-between-promise-and-async-await-in-node-js/)
- [Promise.all vs Multiple Await (Stack Overflow)](https://stackoverflow.com/questions/45285129)

---

## 6. Cross-References

- [TypeScript Standards Reference](./08-typescript-standards-reference.md) — TS conventions
- [Master Coding Guidelines §6](../01-cross-language/15-master-coding-guidelines/00-overview.md) — Error handling

---

*Promise & await patterns — consolidated from pre-code review guides.*
