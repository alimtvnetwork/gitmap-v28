# Subtask: Subagent 3: Nested If Flattening & Function Line Cap Extractions

- parent_plan: 01-coding-guideline-fixes.md
- steps: 101 to 150
- status: pending

## Execution Instructions
- Strictly enforce the <= 15 lines rule for any modified or extracted functions.
- Wrap errors with `apperror.Wrap` and domain context.
- Never use negative booleans or `!isX` inverted logic.
- Place exactly one blank line before `return` statements and after closing `}` braces.

## Granular Steps

101. **nested-if**: `gitmap/archive/extract.go:192` - Nested `if` statement detected at indentation depth 3.
   - **Action**: Flatten with early returns or guard clauses.
102. **nested-if**: `gitmap/archive/extract.go:233` - Nested `if` statement detected at indentation depth 3.
   - **Action**: Flatten with early returns or guard clauses.
103. **nested-if**: `gitmap/archive/extract.go:265` - Nested `if` statement detected at indentation depth 3.
   - **Action**: Flatten with early returns or guard clauses.
104. **nested-if**: `gitmap/archive/extract.go:293` - Nested `if` statement detected at indentation depth 3.
   - **Action**: Flatten with early returns or guard clauses.
105. **nested-if**: `gitmap/archive/source.go:175` - Nested `if` statement detected at indentation depth 3.
   - **Action**: Flatten with early returns or guard clauses.
106. **nested-if**: `gitmap/archive/source.go:195` - Nested `if` statement detected at indentation depth 3.
   - **Action**: Flatten with early returns or guard clauses.
107. **nested-if**: `gitmap/archive/source.go:199` - Nested `if` statement detected at indentation depth 4.
   - **Action**: Flatten with early returns or guard clauses.
108. **nested-if**: `gitmap/cliexit/report.go:183` - Nested `if` statement detected at indentation depth 3.
   - **Action**: Flatten with early returns or guard clauses.
109. **nested-if**: `gitmap/clonefrom/parsecsv.go:55` - Nested `if` statement detected at indentation depth 3.
   - **Action**: Flatten with early returns or guard clauses.
110. **nested-if**: `gitmap/clonefrom/render.go:113` - Nested `if` statement detected at indentation depth 3.
   - **Action**: Flatten with early returns or guard clauses.
111. **nested-if**: `gitmap/clonefrom/summary.go:43` - Nested `if` statement detected at indentation depth 3.
   - **Action**: Flatten with early returns or guard clauses.
112. **nested-if**: `gitmap/clonefrom/summary.go:47` - Nested `if` statement detected at indentation depth 4.
   - **Action**: Flatten with early returns or guard clauses.
113. **nested-if**: `gitmap/clonefrom/summary.go:117` - Nested `if` statement detected at indentation depth 3.
   - **Action**: Flatten with early returns or guard clauses.
114. **nested-if**: `gitmap/clonefrom/summary_terminal.go:51` - Nested `if` statement detected at indentation depth 3.
   - **Action**: Flatten with early returns or guard clauses.
115. **nested-if**: `gitmap/clonefrom/summary_terminal.go:73` - Nested `if` statement detected at indentation depth 3.
   - **Action**: Flatten with early returns or guard clauses.
116. **nested-if**: `gitmap/clonefrom/summary_terminal.go:146` - Nested `if` statement detected at indentation depth 3.
   - **Action**: Flatten with early returns or guard clauses.
117. **nested-if**: `gitmap/clonefrom/validate.go:44` - Nested `if` statement detected at indentation depth 3.
   - **Action**: Flatten with early returns or guard clauses.
118. **nested-if**: `gitmap/clonefrom/validate.go:47` - Nested `if` statement detected at indentation depth 4.
   - **Action**: Flatten with early returns or guard clauses.
119. **nested-if**: `gitmap/clonefrom/validate.go:50` - Nested `if` statement detected at indentation depth 5.
   - **Action**: Flatten with early returns or guard clauses.
120. **nested-if**: `gitmap/clonefrom/validate.go:181` - Nested `if` statement detected at indentation depth 3.
   - **Action**: Flatten with early returns or guard clauses.
121. **nested-if**: `gitmap/clonenext/localstate.go:123` - Nested `if` statement detected at indentation depth 3.
   - **Action**: Flatten with early returns or guard clauses.
122. **nested-if**: `gitmap/clonenext/repodetect.go:46` - Nested `if` statement detected at indentation depth 3.
   - **Action**: Flatten with early returns or guard clauses.
123. **nested-if**: `gitmap/clonenow/clonenow.go:102` - Nested `if` statement detected at indentation depth 3.
   - **Action**: Flatten with early returns or guard clauses.
124. **nested-if**: `gitmap/clonenow/execute.go:97` - Nested `if` statement detected at indentation depth 3.
   - **Action**: Flatten with early returns or guard clauses.
125. **nested-if**: `gitmap/clonenow/execute.go:141` - Nested `if` statement detected at indentation depth 3.
   - **Action**: Flatten with early returns or guard clauses.
126. **long-func**: `gitmap/archive/create.go:108` - Function `CreateArchive` exceeds 15 lines (20 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
127. **long-func**: `gitmap/archive/extract.go:38` - Function `CompactExtract` exceeds 15 lines (16 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
128. **long-func**: `gitmap/archive/extract.go:66` - Function `completeCompactExtract` exceeds 15 lines (20 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
129. **long-func**: `gitmap/archive/extract.go:92` - Function `extractAllIntoDir` exceeds 15 lines (20 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
130. **long-func**: `gitmap/archive/extract.go:125` - Function `extractArchiveEntry` exceeds 15 lines (17 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
131. **long-func**: `gitmap/archive/extract.go:177` - Function `safeJoin` exceeds 15 lines (21 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
132. **long-func**: `gitmap/archive/extract.go:202` - Function `promoteRealRoot` exceeds 15 lines (16 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
133. **long-func**: `gitmap/archive/extract.go:219` - Function `findDeepestRoot` exceeds 15 lines (22 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
134. **long-func**: `gitmap/archive/extract.go:282` - Function `copyDirEntry` exceeds 15 lines (17 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
135. **long-func**: `gitmap/archive/extract.go:327` - Function `archiveBaseName` exceeds 15 lines (16 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
136. **long-func**: `gitmap/archive/list.go:29` - Function `ListEntries` exceeds 15 lines (20 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
137. **long-func**: `gitmap/archive/list.go:50` - Function `extractListEntries` exceeds 15 lines (16 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
138. **long-func**: `gitmap/archive/source.go:96` - Function `ResolveSource` exceeds 15 lines (22 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
139. **long-func**: `gitmap/archive/source.go:131` - Function `resolveHTTP` exceeds 15 lines (21 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
140. **long-func**: `gitmap/archive/source.go:157` - Function `downloadWithAria2c` exceeds 15 lines (24 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
141. **long-func**: `gitmap/archive/source.go:184` - Function `downloadWithHTTP` exceeds 15 lines (23 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
142. **long-func**: `gitmap/archive/source.go:227` - Function `resolveGit` exceeds 15 lines (17 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
143. **long-func**: `gitmap/archive/source.go:249` - Function `AutoDetectSingleArchive` exceeds 15 lines (28 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
144. **long-func**: `gitmap/cliexit/report.go:149` - Function `sortedExtraLines` exceeds 15 lines (16 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
145. **long-func**: `gitmap/cliexit/report.go:168` - Function `writeJSON` exceeds 15 lines (24 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
146. **long-func**: `gitmap/clonefrom/execute_hooks.go:34` - Function `ExecuteWithHooks` exceeds 15 lines (16 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
147. **long-func**: `gitmap/clonefrom/execute_lfs_fix.go:39` - Function `executeLFSFix` exceeds 15 lines (38 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
148. **long-func**: `gitmap/clonefrom/jsonschema_helpers.go:41` - Function `rootSchema` exceeds 15 lines (17 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
149. **long-func**: `gitmap/clonefrom/parse.go:33` - Function `ParseFile` exceeds 15 lines (20 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
150. **long-func**: `gitmap/clonefrom/parse.go:81` - Function `parseJSON` exceeds 15 lines (17 lines).
   - **Action**: Extract helper functions, table-driven dispatch, or guard clauses.
