import esbuild from 'esbuild';
import esserve from '@es-exec/esbuild-plugin-serve';

esbuild.build({
  entryPoints: ['src/server.ts'],
  outdir: 'dist',
  plugins: [esserve()],
  format: 'esm',
  watch: true,
});
