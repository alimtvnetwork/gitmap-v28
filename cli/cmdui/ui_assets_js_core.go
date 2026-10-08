package cmdui

const uiAssetsJSCore = `    function showTab(tabId) {
      document.querySelectorAll('.content-area').forEach(el => el.classList.remove('active'));
      document.querySelectorAll('.nav-btn').forEach(el => el.classList.remove('active'));
      const target = document.getElementById('tab-' + tabId);
      if (target) target.classList.add('active');
      document.getElementById('header-title').innerText = tabId.charAt(0).toUpperCase() + tabId.slice(1);
      event?.target?.classList.add('active');
      if (tabId === 'ssh') refreshNodes();
      if (tabId === 'editor') populateEditorNodes();
    }

`
