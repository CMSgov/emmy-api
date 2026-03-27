import { readFile } from 'node:fs/promises';
import Ajv2020 from 'ajv/dist/2020.js';
import addFormats from 'ajv-formats';
import { glob } from 'glob';

const ajv = new Ajv2020({
  strict: true,
  allErrors: true,
  validateFormats: true
});

addFormats(ajv);

const files = (await glob('schema/**/*.schema.json')).sort();

if (files.length === 0) {
  throw new Error('No schema files found under schema/');
}

const schemas = await Promise.all(
  files.map(async (path) => {
    const text = await readFile(path, 'utf8');
    return {
      path,
      schema: JSON.parse(text)
    };
  })
);

for (const { path, schema } of schemas) {
  console.info(`Loading ${path}`);
  ajv.addSchema(schema);
}

for (const { path, schema } of schemas) {
  console.info(`Validating ${path}`);

  if (typeof schema.$id === 'string' && schema.$id.length > 0) {
    ajv.getSchema(schema.$id);
  } else {
    ajv.compile(schema);
  }
}

console.info(`Validated ${schemas.length} schema files with strict draft 2020-12.`);
