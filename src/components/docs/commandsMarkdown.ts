import { CommandCategory, CommandDef, CommandExample, CommandFlag } from "@/data/commands";

function formatCommandFlags(flags?: CommandFlag[]): string {
  const hasFlags = Boolean(flags && flags.length > 0);
  if (hasFlags === false) return "";

  const lines = flags!.map((item) => `- \`${item.flag}\` — ${item.description}`);

  return `**Flags:**\n${lines.join("\n")}\n\n`;
}

function formatCommandExamples(examples?: CommandExample[]): string {
  const hasExamples = Boolean(examples && examples.length > 0);
  if (hasExamples === false) return "";

  const formatted = examples!.map((item) => {
    const desc = item.description ? `${item.description}:\n` : "";
    return `${desc}\`\`\`bash\n${item.command}\n\`\`\``;
  });

  return `**Examples:**\n${formatted.join("\n\n")}\n\n`;
}

function formatSingleCommand(commandItem: CommandDef): string {
  const aliasText = commandItem.alias ? ` (alias: \`${commandItem.alias}\`)` : "";
  const usageText = commandItem.usage ? `**Usage:**\n\`\`\`\n${commandItem.usage}\n\`\`\`\n\n` : "";
  const flagsText = formatCommandFlags(commandItem.flags);
  const examplesText = formatCommandExamples(commandItem.examples);

  return `### \`${commandItem.name}\`${aliasText}\n\n${commandItem.description}\n\n${usageText}${flagsText}${examplesText}`;
}

function formatCategoryGroup(category: CommandCategory, allCommands: CommandDef[]): string {
  const categoryCommands = allCommands.filter((cmd) => cmd.category === category.key);
  if (categoryCommands.length === 0) return "";

  const iconText = category.icon ? `${category.icon} ` : "";
  const header = `## ${iconText}${category.label}\n\n${category.description}\n\n`;
  const body = categoryCommands.map(formatSingleCommand).join("");

  return `${header}${body}---\n\n`;
}

export function generateCommandsMarkdown(allCommands: CommandDef[], categories: CommandCategory[]): string {
  const intro = `# gitmap Command Reference\n\n> ${allCommands.length} commands organized by category.\n\n`;
  const groups = categories.map((category) => formatCategoryGroup(category, allCommands)).join("");

  return `${intro}${groups}`;
}

export function downloadMarkdownFile(content: string, fileName: string): void {
  const blob = new Blob([content], { type: "text/markdown" });
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement("a");
  anchor.href = url;
  anchor.download = fileName;
  anchor.click();
  URL.revokeObjectURL(url);
}
