Part 1

Output is still wrong for 

gitmap nodes clone https://github.com/alimtvnetwork/awansoft-v10

I can still acess to w3 machine fix it and release minor please and 

PS C:\Users\Alim\awansoft-v10> gitmap pa --ssh

  Enqueuing 'pull-all' across SSH fleet:
    • Remote Node [w3] (192.168.1.12): Offline (skipped, no task enqueued)
    • Remote Node [w1] (192.168.1.3): Online → Enqueued (async)
    • Remote Node [w2] (192.168.1.7): Online → Enqueued (async)
    • Remote Node [w4] (192.168.1.13): Online → Enqueued (async)
    • Remote Node [u1] (192.168.1.22): Offline (skipped, no task enqueued)
    • Remote Node [main] (192.168.1.20): Offline (skipped, no task enqueued)
    • Current Machine [Alim-Desktop (127.0.0.1)]: Running locally (direct execution, not enqueued)

  ▶ Local VM (127.0.0.1 - localhost): 45 pulled (22 active, 23 up-to-date)
...
(Includes failure regarding credential store on Windows nodes)

Additional commands mentioned:
gitmap pull all (pa)
gitmap pull all ssh (pas)
gitmap fix ignore all  [-y]# will fix on all repos to fix gitignroe issues
gitmap fix-ignore-all (fia)  [-y] # will fix on all repos to fix gitignroe issues
gitmap fix ignores all ssh [-y] # will fix on all repos to fix gitignroe issues
gitmap fix-ignores-all-ssh (fias)  [-y] # will fix on all repos to fix gitignroe issues for all nodes, current node will run as is, follow the gitmap pas formula for this and write in the spec as a term for Gitmap PAS formula so that can be reffered by to any AI properly
gitmap commit-push-all-repos (cpar) [-y]# commit all repos if there is any pending
gitmap commit-push-all-repos (cpar) --review(r) # show the pending commits for review first if approved then commits and push
gitmap commit-push-all-repos (cpar) --review --commit-only(co) # show the pending commits for review first if approved then commits and push
gitmap see commit pending
gitmap see git-ignore/ignore/ig issues
gitmap see errors # same as gitmap errors
gitmap see history # same gitmap history

gitmap see errors ssh (ses) # same as gitmap errors, ssh follow gitmap PAS
gitmap pull-all-ssh(pas) # do pull all right now ly/optimize
gitmap ignore/ig add/scan/scan-ssh(ss)/remove/edit/action/ls/help/ui/app/add-group/remove-group/rm-grp/set-default-group/add-grp-to-default(agtd) <name of the group>/apply ./connect-group-with-repo (cgwp) /export/import
gitmap ignore connect-group-with-repo (cgwp) <group-name> <path1>,<alias of the repo>
