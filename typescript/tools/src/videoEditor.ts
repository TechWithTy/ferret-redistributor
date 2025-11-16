import fs from 'node:fs';
import path from 'node:path';
import process from 'node:process';
import JSON5 from 'json5';
import editly, { EditlyOptions } from 'editly';

type CLIOptions = {
  configPath: string;
  output?: string;
  width?: number;
  height?: number;
  fast?: boolean;
};

function parseArgs(argv: string[]): CLIOptions {
  const args = [...argv];
  const opts: Partial<CLIOptions> = {};

  while (args.length > 0) {
    const token = args.shift();
    if (!token) break;
    switch (token) {
      case '--config':
      case '-c':
        opts.configPath = requireValue(token, args);
        break;
      case '--output':
      case '-o':
        opts.output = requireValue(token, args);
        break;
      case '--width':
        opts.width = Number(requireValue(token, args));
        break;
      case '--height':
        opts.height = Number(requireValue(token, args));
        break;
      case '--fast':
        opts.fast = true;
        break;
      default:
        throw new Error(`Unknown flag: ${token}`);
    }
  }

  if (!opts.configPath) {
    throw new Error('Missing --config path to an editly JSON/JSON5 config');
  }
  return opts as CLIOptions;
}

function requireValue(flag: string, args: string[]): string {
  const value = args.shift();
  if (!value) {
    throw new Error(`Flag ${flag} requires a value`);
  }
  return value;
}

function loadConfig(configPath: string): EditlyOptions {
  const absolutePath = path.resolve(configPath);
  const contents = fs.readFileSync(absolutePath, 'utf8');
  if (configPath.endsWith('.json5')) {
    return JSON5.parse(contents);
  }
  return JSON.parse(contents);
}

function applyOverrides(config: EditlyOptions, opts: CLIOptions): EditlyOptions {
  if (opts.output) {
    config.outPath = opts.output;
  }
  if (opts.width || opts.height) {
    config.width = opts.width ?? config.width;
    config.height = opts.height ?? config.height;
  }
  if (opts.fast) {
    config.fast = true;
  }
  return config;
}

async function main() {
  const argv = process.argv.slice(2);
  const opts = parseArgs(argv);
  const config = applyOverrides(loadConfig(opts.configPath), opts);

  console.log('[video-editor] rendering with config:', config.outPath ?? 'preview');
  await editly(config);
  console.log('[video-editor] render finished');
}

main().catch((err) => {
  console.error('[video-editor] failed:', err);
  process.exit(1);
});




