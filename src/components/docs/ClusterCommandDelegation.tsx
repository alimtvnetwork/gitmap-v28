const ClusterCommandDelegation = () => {
  return (
    <div className="mb-8 p-6 bg-muted/30 rounded-lg border border-border">
      <h2 className="text-xl font-bold mb-3 docs-h2">Cluster Command Delegation</h2>
      <p className="text-sm text-muted-foreground mb-4">
        Once <code className="text-primary">gitmap serve</code> and <code className="text-primary">gitmap join</code> establish a cluster of machines, operators use a unified CLI surface to broadcast shell commands, Git operations, or package installations across the network.
      </p>
      <ul className="list-disc list-inside text-sm text-muted-foreground space-y-2 mb-4">
        <li><strong>Target Selectors</strong>: <code className="text-primary">servers-clients</code> targets all nodes including the server, while <code className="text-primary">clients</code> excludes the server.</li>
        <li><strong>Filtering</strong>: Use <code className="text-primary">--except &lt;list&gt;</code> to exclude nodes by IP or ID, or <code className="text-primary">--ip</code> and <code className="text-primary">--id</code> to target exactly.</li>
        <li><strong>Chaining</strong>: You can chain multiple sub-commands sequentially on targeted nodes (e.g., <code className="text-primary">gitmap clients cmd "whoami", ps "Get-Date" --id 1,3</code>).</li>
        <li><strong>Audit Trail</strong>: All executions maintain a persistent history verifiable via <code className="text-primary">gitmap cluster history</code>.</li>
      </ul>
      <div className="mt-4 p-4 border border-red-500/20 bg-red-500/10 rounded-md">
        <h3 className="font-semibold text-red-500 mb-2">Security Note</h3>
        <p className="text-sm text-foreground">
          Client node credentials used for lifecycle commands are securely stored as bcrypt hashes. Passwords are never exported (e.g., via <code className="text-primary">cluster export</code>), and all cluster communication is strictly encrypted over TLS.
        </p>
      </div>
    </div>
  );
};

export default ClusterCommandDelegation;
