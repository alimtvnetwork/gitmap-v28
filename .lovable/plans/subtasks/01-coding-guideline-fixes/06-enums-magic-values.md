# Subtask: Fix Enums and Magic Values
Status: pending

## Steps
1. Edit `.\gitmap\archive\archive.go` at line 32: Type alias Format acting as enum lacks Type suffix -> Rename to FormatType
3. Edit `.\gitmap\archive\create.go` at line 33: Type alias CompressionMode acting as enum lacks Type suffix -> Rename to CompressionModeType
5. Edit `.\gitmap\archive\source.go` at line 30: Type alias SourceKind acting as enum lacks Type suffix -> Rename to SourceKindType
7. Edit `.\gitmap\cliexit\kind.go` at line 37: Type alias Kind acting as enum lacks Type suffix -> Rename to KindType
8. Edit `.\gitmap\cliexit\report.go` at line 50: Type alias OutputMode acting as enum lacks Type suffix -> Rename to OutputModeType
22. Edit `.\gitmap\cloner\audit.go` at line 30: Type alias AuditAction acting as enum lacks Type suffix -> Rename to AuditActionType
37. Edit `.\gitmap\cluster\registry.go` at line 9: Type alias NodeState acting as enum lacks Type suffix -> Rename to NodeStateType
202. Edit `.\gitmap\cmd\commitin\finalize\conflict.go` at line 12: Type alias ConflictDecision acting as enum lacks Type suffix -> Rename to ConflictDecisionType
211. Edit `.\gitmap\committransfer\types.go` at line 17: Type alias Direction acting as enum lacks Type suffix -> Rename to DirectionType
212. Edit `.\gitmap\committransfer\types.go` at line 63: Type alias PreferPolicy acting as enum lacks Type suffix -> Rename to PreferPolicyType
238. Edit `.\gitmap\constants\constants_inject_idempotency.go` at line 44: Type alias InjectKind acting as enum lacks Type suffix -> Rename to InjectKindType
249. Edit `.\gitmap\diff\tree.go` at line 13: Type alias EntryKind acting as enum lacks Type suffix -> Rename to EntryKindType
252. Edit `.\gitmap\errreport\errreport.go` at line 35: Type alias Phase acting as enum lacks Type suffix -> Rename to PhaseType
261. Edit `.\gitmap\ghtoken\ghtoken.go` at line 27: Type alias Source acting as enum lacks Type suffix -> Rename to SourceType
267. Edit `.\gitmap\glyphs\glyphs.go` at line 18: Type alias Mode acting as enum lacks Type suffix -> Rename to ModeType
268. Edit `.\gitmap\logging\jsonlog.go` at line 20: Type alias Level acting as enum lacks Type suffix -> Rename to LevelType
278. Edit `.\gitmap\movemerge\conflict.go` at line 11: Type alias Choice acting as enum lacks Type suffix -> Rename to ChoiceType
279. Edit `.\gitmap\movemerge\diff.go` at line 6: Type alias DiffKind acting as enum lacks Type suffix -> Rename to DiffKindType
281. Edit `.\gitmap\movemerge\types.go` at line 12: Type alias EndpointKind acting as enum lacks Type suffix -> Rename to EndpointKindType
282. Edit `.\gitmap\movemerge\types.go` at line 34: Type alias PreferPolicy acting as enum lacks Type suffix -> Rename to PreferPolicyType
283. Edit `.\gitmap\movemerge\types.go` at line 50: Type alias Direction acting as enum lacks Type suffix -> Rename to DirectionType
290. Edit `.\gitmap\render\preflag.go` at line 22: Type alias PrettyMode acting as enum lacks Type suffix -> Rename to PrettyModeType
294. Edit `.\gitmap\startup\add.go` at line 65: Type alias AddStatus acting as enum lacks Type suffix -> Rename to AddStatusType
308. Edit `.\gitmap\startup\remove.go` at line 33: Type alias RemoveStatus acting as enum lacks Type suffix -> Rename to RemoveStatusType
318. Edit `.\gitmap\startup\winbackend.go` at line 28: Type alias Backend acting as enum lacks Type suffix -> Rename to BackendType
336. Edit `.\gitmap\templates\diff.go` at line 23: Type alias DiffStatus acting as enum lacks Type suffix -> Rename to DiffStatusType
339. Edit `.\gitmap\templates\merge.go` at line 28: Type alias MergeOutcome acting as enum lacks Type suffix -> Rename to MergeOutcomeType
340. Edit `.\gitmap\templates\resolver.go` at line 12: Type alias Source acting as enum lacks Type suffix -> Rename to SourceType
350. Edit `.\gitmap\theme\theme.go` at line 35: Type alias Mode acting as enum lacks Type suffix -> Rename to ModeType
363. Edit `.\scripts\changelog\internal\runner\args.go` at line 11: Type alias Mode acting as enum lacks Type suffix -> Rename to ModeType
10. Edit `.\gitmap\clonefrom\execute_lfs_fix.go` at line 26: Magic string/number in condition: if match := rxSmudgeFatal.FindStringSubmatch(output); len(match) == 2 { -> Extract to named constant or EnumType
11. Edit `.\gitmap\clonefrom\execute_lfs_fix.go` at line 30: Magic string/number in condition: if match := rxSmudgeError.FindStringSubmatch(output); len(match) == 2 { -> Extract to named constant or EnumType
13. Edit `.\gitmap\clonefrom\parse.go` at line 70: Magic string/number in condition: if format == "json" { -> Extract to named constant or EnumType
17. Edit `.\gitmap\clonenext\localstate.go` at line 105: Magic string/number in condition: if len(parts) == 2 { -> Extract to named constant or EnumType
18. Edit `.\gitmap\clonenow\clonenow.go` at line 97: Magic string/number in condition: if mode == "ssh" { -> Extract to named constant or EnumType
20. Edit `.\gitmap\clonenow\parse_schema.go` at line 108: Magic string/number in condition: if name == "httpsUrl" || name == "sshUrl" { -> Extract to named constant or EnumType
21. Edit `.\gitmap\clonenow\parse_schema.go` at line 148: Magic string/number in condition: if name == "httpsUrl" || name == "sshUrl" { -> Extract to named constant or EnumType
24. Edit `.\gitmap\cloner\runners.go` at line 156: Magic string/number in condition: if strings.ToLower(strings.TrimSpace(response)) == "y" { -> Extract to named constant or EnumType
26. Edit `.\gitmap\cluster\exec_git_test.go` at line 46: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
29. Edit `.\gitmap\cluster\exec_lifecycle.go` at line 62: Magic string/number in condition: if node.NodeRole == "server" || node.IsServer { -> Extract to named constant or EnumType
31. Edit `.\gitmap\cluster\exec_ps_test.go` at line 25: Magic string/number in condition: if file == "pwsh" { -> Extract to named constant or EnumType
32. Edit `.\gitmap\cluster\exec_ps_test.go` at line 67: Magic string/number in condition: if file == "pwsh" { -> Extract to named constant or EnumType
33. Edit `.\gitmap\cluster\exec_ps_test.go` at line 70: Magic string/number in condition: if file == "powershell" { -> Extract to named constant or EnumType
34. Edit `.\gitmap\cluster\exec_ps_test.go` at line 90: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
35. Edit `.\gitmap\cluster\exec_ps_test.go` at line 102: Magic string/number in condition: if file == "pwsh" { -> Extract to named constant or EnumType
36. Edit `.\gitmap\cluster\exec_ps_test.go` at line 122: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
39. Edit `.\gitmap\cmd\amend.go` at line 136: Magic string/number in condition: if f.commitHash == "HEAD" { -> Extract to named constant or EnumType
41. Edit `.\gitmap\cmd\amendexec.go` at line 106: Magic string/number in condition: if f.commitHash == "HEAD" { -> Extract to named constant or EnumType
43. Edit `.\gitmap\cmd\chromeprofile_import_csv.go` at line 65: Magic string/number in condition: if key == "name" { -> Extract to named constant or EnumType
45. Edit `.\gitmap\cmd\cliexit_helpers_test.go` at line 89: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
47. Edit `.\gitmap\cmd\cliexit_scan_test.go` at line 50: Magic string/number in condition: if tc.name == "failure_missing_dir" { -> Extract to named constant or EnumType
61. Edit `.\gitmap\cmd\cluster_ops.go` at line 122: Magic string/number in condition: if format == "csv" { -> Extract to named constant or EnumType
62. Edit `.\gitmap\cmd\cluster_ops.go` at line 325: Magic string/number in condition: if strings.ToLower(n.Status) == "offline" || strings.ToLower(n.Status) == "unreachable" { -> Extract to named constant or EnumType
64. Edit `.\gitmap\cmd\codingguidelines.go` at line 43: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
65. Edit `.\gitmap\cmd\codingguidelines_test.go` at line 19: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
66. Edit `.\gitmap\cmd\codingguidelines_test.go` at line 54: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
67. Edit `.\gitmap\cmd\codingguidelines_test.go` at line 74: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
71. Edit `.\gitmap\cmd\doctorchecks.go` at line 213: Magic string/number in condition: if status == "Valid" { -> Extract to named constant or EnumType
72. Edit `.\gitmap\cmd\doctorchecks.go` at line 220: Magic string/number in condition: if status == "NotSigned" { -> Extract to named constant or EnumType
73. Edit `.\gitmap\cmd\doctordupbin.go` at line 31: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
74. Edit `.\gitmap\cmd\doctordupbin.go` at line 105: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
77. Edit `.\gitmap\cmd\envplatform_unix.go` at line 62: Magic string/number in condition: if shell == "zsh" { -> Extract to named constant or EnumType
78. Edit `.\gitmap\cmd\env_unit_test.go` at line 59: Magic string/number in condition: if v.Name == "B" { -> Extract to named constant or EnumType
79. Edit `.\gitmap\cmd\escapecwd_test.go` at line 76: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
80. Edit `.\gitmap\cmd\expandhome_test.go` at line 56: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
90. Edit `.\gitmap\cmd\installcleancode.go` at line 70: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
93. Edit `.\gitmap\cmd\installctx_harness_test.go` at line 109: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
94. Edit `.\gitmap\cmd\installdetect.go` at line 21: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
95. Edit `.\gitmap\cmd\installdetect.go` at line 24: Magic string/number in condition: if runtime.GOOS == "darwin" { -> Extract to named constant or EnumType
96. Edit `.\gitmap\cmd\installscripts.go` at line 94: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
105. Edit `.\gitmap\cmd\list.go` at line 57: Magic string/number in condition: if lower == "groups" { -> Extract to named constant or EnumType
106. Edit `.\gitmap\cmd\list.go` at line 68: Magic string/number in condition: if lower == "groups" { -> Extract to named constant or EnumType
109. Edit `.\gitmap\cmd\llmdocs.go` at line 55: Magic string/number in condition: if *format == "json" { -> Extract to named constant or EnumType
110. Edit `.\gitmap\cmd\llmdocs.go` at line 111: Magic string/number in condition: if format == "json" { -> Extract to named constant or EnumType
114. Edit `.\gitmap\cmd\llmdocsjson_contract_test.go` at line 100: Magic string/number in condition: if len(got) == 4 && got[3] != "example" { -> Extract to named constant or EnumType
116. Edit `.\gitmap\cmd\move.go` at line 18: Magic string/number in condition: if len(positional) == 2 { -> Extract to named constant or EnumType
121. Edit `.\gitmap\cmd\prune.go` at line 59: Magic string/number in condition: if answer == "y" || answer == "Y" { -> Extract to named constant or EnumType
126. Edit `.\gitmap\cmd\replace.go` at line 67: Magic string/number in condition: if len(positional) == 2 { -> Extract to named constant or EnumType
144. Edit `.\gitmap\cmd\scan_export_clonefrom_integration_test.go` at line 110: Magic string/number in condition: if format == "json" { -> Extract to named constant or EnumType
146. Edit `.\gitmap\cmd\selfinstall.go` at line 248: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
147. Edit `.\gitmap\cmd\selfinstall.go` at line 300: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
148. Edit `.\gitmap\cmd\selfinstall.go` at line 309: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
150. Edit `.\gitmap\cmd\selfuninstallparts.go` at line 65: Magic string/number in condition: if lower == "gitmap" || lower == "gitmap.exe" { -> Extract to named constant or EnumType
154. Edit `.\gitmap\cmd\sshgen.go` at line 130: Magic string/number in condition: if input == "R" { -> Extract to named constant or EnumType
155. Edit `.\gitmap\cmd\sshgen.go` at line 136: Magic string/number in condition: if input == "N" { -> Extract to named constant or EnumType
156. Edit `.\gitmap\cmd\sshgen.go` at line 218: Magic string/number in condition: if action == "update" { -> Extract to named constant or EnumType
157. Edit `.\gitmap\cmd\sshgen.go` at line 221: Magic string/number in condition: if action == "save" { -> Extract to named constant or EnumType
164. Edit `.\gitmap\cmd\templatesdiff.go` at line 140: Magic string/number in condition: if kind == "attributes" { -> Extract to named constant or EnumType
166. Edit `.\gitmap\cmd\update.go` at line 138: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
170. Edit `.\gitmap\cmd\updateremoteinstall.go` at line 108: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
171. Edit `.\gitmap\cmd\updateremoteinstall.go` at line 129: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
172. Edit `.\gitmap\cmd\updateremoteinstall.go` at line 136: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
173. Edit `.\gitmap\cmd\updateremoteinstall.go` at line 155: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
174. Edit `.\gitmap\cmd\updateremoteinstall_test.go` at line 26: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
175. Edit `.\gitmap\cmd\updatescript.go` at line 17: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
178. Edit `.\gitmap\cmd\visibilitybulkprompt.go` at line 98: Magic string/number in condition: if tok == "y" || tok == "yes" { -> Extract to named constant or EnumType
179. Edit `.\gitmap\cmd\vscodepmsync_dedupe_test.go` at line 38: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
180. Edit `.\gitmap\cmd\vscodepmsync_dedupe_test.go` at line 70: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
181. Edit `.\gitmap\cmd\vscodepmsync_dedupe_test.go` at line 106: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
182. Edit `.\gitmap\cmd\vscodepmsync_dedupe_test.go` at line 136: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
184. Edit `.\gitmap\cmd\vscodepmsync_mode_test.go` at line 34: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
185. Edit `.\gitmap\cmd\vscodepmsync_mode_test.go` at line 60: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
186. Edit `.\gitmap\cmd\vscodepmsync_mode_test.go` at line 91: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
187. Edit `.\gitmap\cmd\vscodepmsync_mode_test.go` at line 123: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
188. Edit `.\gitmap\cmd\vscodepmsync_pathtag_test.go` at line 66: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
189. Edit `.\gitmap\cmd\vscodepmsync_pathtag_test.go` at line 92: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
190. Edit `.\gitmap\cmd\vscodepmsync_pathtag_test.go` at line 126: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
191. Edit `.\gitmap\cmd\vscodepmsync_test.go` at line 83: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
192. Edit `.\gitmap\cmd\vscodepmsync_test.go` at line 111: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
194. Edit `.\gitmap\cmd\watchformat.go` at line 56: Magic string/number in condition: if snap.Status == "error" { -> Extract to named constant or EnumType
195. Edit `.\gitmap\cmd\watchformat.go` at line 77: Magic string/number in condition: if status == "dirty" { -> Extract to named constant or EnumType
196. Edit `.\gitmap\cmd\watchops.go` at line 95: Magic string/number in condition: if snap.Status == "dirty" { -> Extract to named constant or EnumType
197. Edit `.\gitmap\cmd\whoami.go` at line 58: Magic string/number in condition: if strings.HasSuffix(name, ".pub") || name == "known_hosts" || name == "config" || strings.HasPrefix(name, "known_hosts") { -> Extract to named constant or EnumType
198. Edit `.\gitmap\cmd\whoami.go` at line 213: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
203. Edit `.\gitmap\cmd\commitin\message\message_test.go` at line 109: Magic string/number in condition: if out == "Refined" { -> Extract to named constant or EnumType
208. Edit `.\gitmap\cmd\commitin\workspace\workspace_test.go` at line 189: Magic string/number in condition: if sub == "clone" && len(args) == 2 { -> Extract to named constant or EnumType
210. Edit `.\gitmap\committransfer\replay.go` at line 144: Magic string/number in condition: if first == "node_modules" && !opts.IncludeNodeMod { -> Extract to named constant or EnumType
213. Edit `.\gitmap\completion\bash.go` at line 56: Magic string/number in condition: if [[ ${COMP_CWORD} -ge 3 ]] && [[ "$sub" == "add" || "$sub" == "show" || "$sub" == "delete" || "$sub" == "remove" || "$sub" == "rename" ]]; then -> Extract to named constant or EnumType
215. Edit `.\gitmap\completion\detect.go` at line 11: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
216. Edit `.\gitmap\completion\dynamic.go` at line 120: Magic string/number in condition: if name == "Default" || strings.HasPrefix(name, "Profile ") { -> Extract to named constant or EnumType
217. Edit `.\gitmap\completion\install.go` at line 124: Magic string/number in condition: if goos == "windows" { -> Extract to named constant or EnumType
226. Edit `.\gitmap\completion\zsh.go` at line 10: Magic string/number in condition: if (( CURRENT == 2 )); then -> Extract to named constant or EnumType
227. Edit `.\gitmap\completion\zsh.go` at line 73: Magic string/number in condition: if (( CURRENT >= 4 )) && [[ "${words[3]}" == "add" || "${words[3]}" == "show" || "${words[3]}" == "delete" || "${words[3]}" == "remove" || "${words[3]}" == "rename" ]]; then -> Extract to named constant or EnumType
229. Edit `.\gitmap\config\config_test.go` at line 14: Magic string/number in condition: if cfg.DefaultMode == "https" { -> Extract to named constant or EnumType
230. Edit `.\gitmap\config\config_test.go` at line 17: Magic string/number in condition: if cfg.DefaultOutput == "terminal" { -> Extract to named constant or EnumType
231. Edit `.\gitmap\config\config_test.go` at line 31: Magic string/number in condition: if cfg.DefaultMode == "https" { -> Extract to named constant or EnumType
232. Edit `.\gitmap\config\config_test.go` at line 51: Magic string/number in condition: if cfg.DefaultMode == "ssh" { -> Extract to named constant or EnumType
233. Edit `.\gitmap\config\config_test.go` at line 64: Magic string/number in condition: if merged.DefaultMode == "ssh" { -> Extract to named constant or EnumType
234. Edit `.\gitmap\config\config_test.go` at line 67: Magic string/number in condition: if merged.DefaultOutput == "json" { -> Extract to named constant or EnumType
235. Edit `.\gitmap\config\config_test.go` at line 84: Magic string/number in condition: if merged.DefaultMode == "ssh" { -> Extract to named constant or EnumType
248. Edit `.\gitmap\desktop\resolve.go` at line 48: Magic string/number in condition: if runtime.GOOS == "darwin" { -> Extract to named constant or EnumType
250. Edit `.\gitmap\diff\tree.go` at line 92: Magic string/number in condition: if !opts.IncludeNodeModules && base == "node_modules" { -> Extract to named constant or EnumType
253. Edit `.\gitmap\formatter\clonescript.go` at line 49: Magic string/number in condition: if r.IdentifiedTransport == "ssh" { -> Extract to named constant or EnumType
257. Edit `.\gitmap\formatter\formatter_test.go` at line 87: Magic string/number in condition: if len(parsed) == 2 { -> Extract to named constant or EnumType
258. Edit `.\gitmap\formatter\formatter_test.go` at line 107: Magic string/number in condition: if len(parsed) == 2 { -> Extract to named constant or EnumType
263. Edit `.\gitmap\gitutil\gitutil.go` at line 279: Magic string/number in condition: if len(parts) == 2 { -> Extract to named constant or EnumType
264. Edit `.\gitmap\glyphs\autodetect_windows.go` at line 24: Magic string/number in condition: if os.Getenv("TERM_PROGRAM") == "vscode" { -> Extract to named constant or EnumType
265. Edit `.\gitmap\glyphs\autodetect_windows.go` at line 27: Magic string/number in condition: if os.Getenv("ConEmuANSI") == "ON" { -> Extract to named constant or EnumType
270. Edit `.\gitmap\macro\execute.go` at line 46: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
271. Edit `.\gitmap\macro\record.go` at line 43: Magic string/number in condition: if cmdText == "stop" || cmdText == "exit" || cmdText == "quit" { -> Extract to named constant or EnumType
272. Edit `.\gitmap\macro\record.go` at line 69: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
274. Edit `.\gitmap\mapper\mapper.go` at line 166: Magic string/number in condition: if len(parts) == 2 { -> Extract to named constant or EnumType
275. Edit `.\gitmap\mapper\mapper.go` at line 177: Magic string/number in condition: if len(parts) == 2 { -> Extract to named constant or EnumType
276. Edit `.\gitmap\mapper\mapper_relroot_test.go` at line 73: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
277. Edit `.\gitmap\mapper\mapper_test.go` at line 61: Magic string/number in condition: if name == "unknown" { -> Extract to named constant or EnumType
284. Edit `.\gitmap\movemerge\walk.go` at line 55: Magic string/number in condition: if base == "node_modules" { -> Extract to named constant or EnumType
286. Edit `.\gitmap\release\assetsbuild.go` at line 48: Magic string/number in condition: if target.GOOS == "windows" { -> Extract to named constant or EnumType
287. Edit `.\gitmap\release\autocommit.go` at line 176: Magic string/number in condition: if answer == "y" || answer == "yes" { -> Extract to named constant or EnumType
295. Edit `.\gitmap\startup\add.go` at line 115: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
296. Edit `.\gitmap\startup\add.go` at line 141: Magic string/number in condition: if runtime.GOOS == "darwin" { -> Extract to named constant or EnumType
297. Edit `.\gitmap\startup\add.go` at line 189: Magic string/number in condition: if runtime.GOOS == "darwin" { -> Extract to named constant or EnumType
299. Edit `.\gitmap\startup\add_test.go` at line 128: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
300. Edit `.\gitmap\startup\desktop.go` at line 33: Magic string/number in condition: if runtime.GOOS == "darwin" { -> Extract to named constant or EnumType
301. Edit `.\gitmap\startup\lifecycle_integration_test.go` at line 43: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
302. Edit `.\gitmap\startup\lifecycle_integration_test.go` at line 46: Magic string/number in condition: if runtime.GOOS == "darwin" { -> Extract to named constant or EnumType
303. Edit `.\gitmap\startup\lifecycle_integration_test.go` at line 58: Magic string/number in condition: if runtime.GOOS == "darwin" { -> Extract to named constant or EnumType
304. Edit `.\gitmap\startup\plist.go` at line 143: Magic string/number in condition: if s.pendingKey == "ProgramArguments" { -> Extract to named constant or EnumType
305. Edit `.\gitmap\startup\plist.go` at line 153: Magic string/number in condition: if s.pendingKey == "Program" { -> Extract to named constant or EnumType
306. Edit `.\gitmap\startup\plistxml.go` at line 47: Magic string/number in condition: if start, ok := tok.(xml.StartElement); ok && start.Name.Local == "string" { -> Extract to named constant or EnumType
307. Edit `.\gitmap\startup\plistxml.go` at line 51: Magic string/number in condition: if end, ok := tok.(xml.EndElement); ok && end.Name.Local == "array" { -> Extract to named constant or EnumType
309. Edit `.\gitmap\startup\remove.go` at line 88: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
310. Edit `.\gitmap\startup\remove.go` at line 175: Magic string/number in condition: if runtime.GOOS == "darwin" { -> Extract to named constant or EnumType
311. Edit `.\gitmap\startup\remove.go` at line 189: Magic string/number in condition: if runtime.GOOS == "darwin" { -> Extract to named constant or EnumType
312. Edit `.\gitmap\startup\startup.go` at line 71: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
313. Edit `.\gitmap\startup\startup.go` at line 75: Magic string/number in condition: if runtime.GOOS == "darwin" { -> Extract to named constant or EnumType
314. Edit `.\gitmap\startup\startup.go` at line 111: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
315. Edit `.\gitmap\startup\startup_test.go` at line 19: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
316. Edit `.\gitmap\startup\startup_test.go` at line 22: Magic string/number in condition: if runtime.GOOS == "darwin" { -> Extract to named constant or EnumType
317. Edit `.\gitmap\startup\startup_test.go` at line 84: Magic string/number in condition: if runtime.GOOS == "windows" || runtime.GOOS == "darwin" { -> Extract to named constant or EnumType
319. Edit `.\gitmap\startup\winbackend.go` at line 96: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
322. Edit `.\gitmap\store\downloader_seed.go` at line 85: Magic string/number in condition: if arg == "version" || arg == "--version" || arg == "-v" { -> Extract to named constant or EnumType
323. Edit `.\gitmap\store\migrateids.go` at line 54: Magic string/number in condition: if name == "Id" && colType == "TEXT" { -> Extract to named constant or EnumType
337. Edit `.\gitmap\templates\list_test.go` at line 20: Magic string/number in condition: if e.Kind == kindIgnore && e.Lang == "go" { -> Extract to named constant or EnumType
338. Edit `.\gitmap\templates\list_test.go` at line 23: Magic string/number in condition: if e.Kind == kindIgnore && e.Lang == "go" && e.Source != SourceEmbed { -> Extract to named constant or EnumType
341. Edit `.\gitmap\tests\cmd_test\amend_test.go` at line 78: Magic string/number in condition: if commitHash == "HEAD" { -> Extract to named constant or EnumType
347. Edit `.\gitmap\tests\fixrepo_test\gofmt_e2e_test.go` at line 100: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
349. Edit `.\gitmap\tests\store_test\location_test.go` at line 37: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
357. Edit `.\gitmap\vscodepm\io.go` at line 73: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
358. Edit `.\gitmap\vscodepm\path_test.go` at line 128: Magic string/number in condition: if runtime.GOOS == "darwin" { -> Extract to named constant or EnumType
359. Edit `.\gitmap\vscodepm\path_test.go` at line 139: Magic string/number in condition: if runtime.GOOS == "windows" { -> Extract to named constant or EnumType
360. Edit `.\remotion-demo\src\Terminal.tsx` at line 34: Magic string/number in condition: if (ln.kind === "prompt") { -> Extract to named constant or EnumType
361. Edit `.\remotion-demo\src\Terminal.tsx` at line 164: Magic string/number in condition: if (ln.kind === "blank") { -> Extract to named constant or EnumType
362. Edit `.\remotion-demo\src\Terminal.tsx` at line 168: Magic string/number in condition: if (ln.kind === "prompt") { -> Extract to named constant or EnumType
371. Edit `.\src\components\docs\DocsTooltip.tsx` at line 42: Magic string/number in condition: if (typeof label === "string") return label; -> Extract to named constant or EnumType
377. Edit `.\src\components\docs\TabOrderMap.tsx` at line 51: Magic string/number in condition: if (el.getAttribute("aria-hidden") === "true") return false; -> Extract to named constant or EnumType
378. Edit `.\src\components\docs\TabOrderMap.tsx` at line 63: Magic string/number in condition: if (s.display === "none") return false; -> Extract to named constant or EnumType
379. Edit `.\src\components\docs\TabOrderMap.tsx` at line 64: Magic string/number in condition: if (s.visibility === "hidden" || s.visibility === "collapse") return false; -> Extract to named constant or EnumType
380. Edit `.\src\components\docs\TabOrderMap.tsx` at line 66: Magic string/number in condition: if (s.pointerEvents === "none" && node === el) { -> Extract to named constant or EnumType
390. Edit `.\src\components\ui\carousel.tsx` at line 77: Magic string/number in condition: if (event.key === "ArrowLeft") { -> Extract to named constant or EnumType
391. Edit `.\src\components\ui\chart.tsx` at line 285: Magic string/number in condition: if (typeof payload !== "object" || payload === null) { -> Extract to named constant or EnumType
392. Edit `.\src\components\ui\chart.tsx` at line 296: Magic string/number in condition: if (key in payload && typeof payload[key as keyof typeof payload] === "string") { -> Extract to named constant or EnumType
396. Edit `.\src\components\ui\sidebar.tsx` at line 484: Magic string/number in condition: if (typeof tooltip === "string") { -> Extract to named constant or EnumType
400. Edit `.\src\hooks\useTheme.ts` at line 49: Magic string/number in condition: if (next === "light" || next === "dark") setThemeState(next); -> Extract to named constant or EnumType
401. Edit `.\src\hooks\useTheme.ts` at line 54: Magic string/number in condition: if (next === "system" || next === "user") setSourceState(next); -> Extract to named constant or EnumType
402. Edit `.\src\hooks\useTheme.ts` at line 59: Magic string/number in condition: if (event.newValue === "light" || event.newValue === "dark") { -> Extract to named constant or EnumType
403. Edit `.\src\lib\clipboard.ts` at line 32: Magic string/number in condition: if (typeof document === "undefined") { -> Extract to named constant or EnumType
404. Edit `.\src\lib\theme.ts` at line 13: Magic string/number in condition: if (typeof document === "undefined") return ThemeType.Dark; -> Extract to named constant or EnumType
405. Edit `.\src\lib\theme.ts` at line 18: Magic string/number in condition: if (res.isSuccess && (res.data === "light" || res.data === "dark")) { -> Extract to named constant or EnumType
406. Edit `.\src\lib\theme.ts` at line 26: Magic string/number in condition: if (typeof document === "undefined") return; -> Extract to named constant or EnumType
440. Edit `.\src\pages\SpecIndex.tsx` at line 36: Magic string/number in condition: if (e.key === "Escape" && document.activeElement === inputRef.current) { -> Extract to named constant or EnumType
451. Edit `.\src\types\helpJson.ts` at line 30: Magic string/number in condition: if (!value || typeof value !== "object") return false; -> Extract to named constant or EnumType
452. Edit `.\src\types\helpJson.ts` at line 32: Magic string/number in condition: if (typeof v.version !== "string") return false; -> Extract to named constant or EnumType
453. Edit `.\src\types\helpJson.ts` at line 33: Magic string/number in condition: if (typeof v.count !== "number") return false; -> Extract to named constant or EnumType