package cmdui

const uiAssetsJSPromptsHelp = `    let allPrompts = [], activePromptCategory = 'all';
    async function loadPromptsList() {
      const container = document.getElementById('prompts-list');
      if (!container) return;
      try {
        const res = await fetch('/api/prompts/instances?limit=50&include_convs=true');
        const data = await res.json();
        allPrompts = [];
        if (data && data.instances) {
          data.instances.forEach(function (inst) {
            const instName = inst.instanceName || inst.instanceId || 'instance';
            (inst.running || []).forEach(function (p) {
              allPrompts.push({
                title: p.title || p.taskDescription || 'Running Task',
                category: 'active',
                instance: instName,
                status: 'running',
                text: p.promptText || p.taskDescription || ''
              });
            });
            (inst.queued || []).forEach(function (q) {
              allPrompts.push({
                title: q.title || q.prompt || 'Queued Instruction',
                category: 'active',
                instance: instName,
                status: 'queued',
                text: q.prompt || q.text || ''
              });
            });
            (inst.recentConvs || []).forEach(function (c) {
              allPrompts.push({
                title: c.title || ('Conversation ' + (c.conversationId || '').slice(0, 8)),
                category: c.isActive ? 'active' : 'completed',
                instance: instName,
                status: c.isActive ? 'active' : 'idle',
                text: c.lastPromptText || ''
              });
            });
          });
        }
        renderPromptsList();
      } catch (e) {
        container.innerHTML = '<div style="color:var(--error);font-size:var(--fs-sm)">Failed to load prompts: ' + escapeHTML(e.message) + '</div>';
      }
    }
    function setPromptCategory(cat, btn) {
      activePromptCategory = cat;
      document.querySelectorAll('#prompt-category-pills .chip').forEach(function (p) { p.classList.toggle('active', p === btn); });
      renderPromptsList();
    }
    function filterPromptsUI(q) { renderPromptsList(q); }
    function renderPromptsList(queryOverride) {
      const container = document.getElementById('prompts-list');
      if (!container) return;
      const searchEl = document.getElementById('prompt-search');
      const q = (queryOverride !== undefined ? queryOverride : (searchEl ? searchEl.value : '')).toLowerCase().trim();
      let filtered = allPrompts;
      if (activePromptCategory !== 'all') filtered = filtered.filter(p => p.category === activePromptCategory);
      if (q) filtered = filtered.filter(p => (p.title && p.title.toLowerCase().includes(q)) || (p.text && p.text.toLowerCase().includes(q)));
      if (!filtered.length) {
        container.innerHTML = '<div style="color:var(--muted);font-size:var(--fs-sm);padding:var(--sp-2)">No prompts found matching filter.</div>';
        return;
      }
      function rankPromptRow(s) {
        const st = (s || '').toLowerCase();
        if (st === 'running') return 0;
        if (st === 'queued' || st === 'enqueued') return 1;
        if (st === 'active') return 2;
        return 3;
      }
      const promptGroups = [];
      const promptGroupSeen = {};
      filtered.forEach(function (p) {
        const instName = p.instance || 'instance';
        if (promptGroupSeen[instName] === undefined) {
          promptGroupSeen[instName] = promptGroups.length;
          promptGroups.push({ name: instName, rows: [] });
        }
        promptGroups[promptGroupSeen[instName]].rows.push(p);
      });
      container.innerHTML = '';
      let shownPromptCount = 0;
      promptGroups.forEach(function (g) {
        if (shownPromptCount >= 30) return;
        g.rows.sort(function (a, b) { return rankPromptRow(a.status) - rankPromptRow(b.status); });
        const heading = document.createElement('div');
        heading.style.cssText = 'font-size:var(--fs-sm);font-weight:700;color:var(--text);margin-top:var(--sp-2);';
        heading.textContent = 'From instance: ' + g.name;
        container.appendChild(heading);
        g.rows.forEach(function (p) {
          if (shownPromptCount >= 30) return;
          shownPromptCount += 1;
          const item = document.createElement('div');
          item.style.cssText = 'padding:var(--sp-3);background:var(--raised);border:1px solid var(--border);border-radius:var(--radius-sm);cursor:pointer;';
          const stLower = (p.status || '').toLowerCase();
          let badgeClass = 'badge-offline';
          let badgeText = (p.status || 'IDLE').toUpperCase();
          let badgeStyle = '';
          if (stLower === 'running') {
            badgeClass = 'badge-online';
            badgeText = 'RUNNING';
            badgeStyle = 'background:rgba(210,153,34,0.18);color:#d29922;border:1px solid #d29922;font-weight:700;';
          } else if (stLower === 'queued' || stLower === 'enqueued') {
            badgeClass = 'badge-queued';
            badgeText = stLower === 'enqueued' ? 'ENQUEUED' : 'QUEUED';
            badgeStyle = 'background:rgba(88,166,255,0.18);color:#58a6ff;border:1px solid #58a6ff;font-weight:700;';
          } else if (stLower === 'active') {
            badgeClass = 'badge-online';
            badgeText = 'ACTIVE';
          }
          item.innerHTML = '<div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:4px;">' +
            '<b>' + escapeHTML(p.title) + '</b>' +
            '<span class="badge ' + badgeClass + '"' + (badgeStyle ? ' style="' + badgeStyle + '"' : '') + '>' + escapeHTML(badgeText) + '</span>' +
            '</div>' +
            '<div style="font-size:var(--fs-xs);color:var(--muted);margin-bottom:4px;">Source: ' + escapeHTML(p.instance || g.name) + '</div>' +
            '<div style="font-size:var(--fs-sm);color:var(--text);font-family:var(--font-mono);max-height:48px;overflow:hidden;text-overflow:ellipsis;">' + escapeHTML((p.text || '').slice(0, 160)) + '</div>';
          item.addEventListener('click', function () {
            const area = document.getElementById('prompt-content');
            if (area) { area.value = p.text || p.title; area.scrollIntoView({ behavior: 'smooth' }); }
            showToast('Loaded prompt template into scratchpad', 'info');
          });
          container.appendChild(item);
        });
      });
    }
    function copyPromptToClipboard() {
      const area = document.getElementById('prompt-content');
      if (!area || !area.value.trim()) { showToast('No prompt content to copy', 'error'); return; }
      navigator.clipboard.writeText(area.value);
      showToast('Prompt copied to clipboard!', 'success');
    }

    let allHelpTopics = [], activeHelpTopic = '';
    async function loadHelpTopics() {
      const listEl = document.getElementById('help-topics-list');
      if (!listEl) return;
      try {
        const res = await fetch('/api/help/topics');
        const data = await res.json();
        const topics = Array.isArray(data) ? data : (data.topics || []);
        allHelpTopics = topics;
        const countEl = document.getElementById('help-topics-count');
        if (countEl) countEl.innerText = topics.length + ' Topics';
        renderHelpTopicsList();
        if (topics.length > 0 && !activeHelpTopic) {
          const first = topics[0];
          loadHelpDoc(typeof first === 'string' ? first : (first.topic || first.name));
        }
      } catch (e) {
        listEl.innerHTML = '<div style="padding:var(--sp-2);color:var(--error);font-size:var(--fs-sm)">Failed to load topics: ' + escapeHTML(e.message) + '</div>';
      }
    }
    function filterHelpTopicsUI(q) { renderHelpTopicsList(q); }
    function renderHelpTopicsList(queryOverride) {
      const listEl = document.getElementById('help-topics-list');
      if (!listEl) return;
      const searchEl = document.getElementById('help-search');
      const q = (queryOverride !== undefined ? queryOverride : (searchEl ? searchEl.value : '')).toLowerCase().trim();
      let filtered = allHelpTopics;
      if (q) filtered = filtered.filter(t => {
        const n = (typeof t === 'string' ? t : (t.name || t.topic || '')).toLowerCase();
        const s = (typeof t === 'string' ? '' : (t.summary || '')).toLowerCase();
        return n.includes(q) || s.includes(q);
      });
      const visCount = document.getElementById('help-visible-count');
      if (visCount) visCount.innerText = filtered.length;
      if (!filtered.length) {
        listEl.innerHTML = '<div style="padding:var(--sp-2);color:var(--muted);font-size:var(--fs-xs)">No matching topics.</div>';
        return;
      }
      listEl.innerHTML = '';
      filtered.forEach(function (t) {
        const name = typeof t === 'string' ? t : (t.name || t.topic || '');
        const summary = typeof t === 'string' ? '' : (t.summary || '');
        const btn = document.createElement('button');
        btn.type = 'button';
        btn.className = 'nav-btn' + (name === activeHelpTopic ? ' active' : '');
        btn.style.cssText = 'padding:6px 10px;font-size:var(--fs-xs);text-align:left;border-radius:var(--radius-sm);margin-bottom:2px;';
        btn.innerHTML = '<b style="display:block;color:var(--primary)">' + escapeHTML(name) + '</b>' +
          (summary ? '<span style="display:block;color:var(--muted);font-size:0.7rem;line-height:1.2;margin-top:2px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">' + escapeHTML(summary) + '</span>' : '');
        btn.addEventListener('click', function () {
          loadHelpDoc(name);
          document.querySelectorAll('#help-topics-list .nav-btn').forEach(function (b) { b.classList.toggle('active', b === btn); });
        });
        listEl.appendChild(btn);
      });
    }
    async function loadHelpDoc(topic) {
      if (!topic) return;
      activeHelpTopic = topic;
      const titleEl = document.getElementById('help-active-topic-title');
      if (titleEl) titleEl.innerText = 'gitmap ' + topic;
      const viewer = document.getElementById('help-doc-viewer');
      if (viewer) viewer.innerText = 'Loading manual for ' + topic + '...';
      try {
        const res = await fetch('/api/help/doc?topic=' + encodeURIComponent(topic));
        const hasManual = res.ok;
        if (hasManual && viewer) {
          viewer.innerText = await res.text();
          return;
        }
        if (viewer) viewer.innerText = 'The topic ' + topic + ' has no embedded manual yet.';
      } catch (e) {
        console.error('loadHelpDoc', topic, e);
        if (viewer) viewer.innerText = 'The topic ' + topic + ' has no embedded manual yet.';
      }
    }
    function copyCurrentHelpDoc() {
      const viewer = document.getElementById('help-doc-viewer');
      if (!viewer || !viewer.innerText.trim()) { showToast('No document content to copy', 'error'); return; }
      navigator.clipboard.writeText(viewer.innerText);
      showToast('Manual markdown copied to clipboard!', 'success');
    }
    function escapeHTML(str) {
      if (!str) return '';
      return String(str).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
    }

    (function syncDocTitle() {
      if (typeof showTab !== 'function') return;
      const orig = showTab;
      showTab = function (tabId) {
        orig(tabId);
        try { const item = NAV_ITEMS.find(function (n) { return n.id === tabId; }); document.title = 'GitMap — ' + (item ? item.label : tabId); } catch (e) {}
        if (tabId === 'terminal') {
          const termTab = document.getElementById('tab-terminal');
          if (termTab) {
            document.querySelectorAll('.content-area').forEach(function (el) { el.classList.remove('active'); });
            termTab.classList.add('active');
          }
          const termIn = document.getElementById('term-input');
          if (termIn) setTimeout(function () { termIn.focus(); }, 60);
        }
        if (tabId === 'prompts') loadPromptsList();
        if (tabId === 'help') loadHelpTopics();
      };
    })();
    window.addEventListener('load', function () {
      renderNav();
      wireTheme();
      wireSettingsSubTabs();
      const p = window.location.pathname.replace(/^\//, '');
      if (p === 'terminal') {
        showTab('terminal');
      }
      document.title = 'GitMap — Settings';
    });
`
