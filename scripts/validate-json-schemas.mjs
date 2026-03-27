import { readFile } from "node:fs/promises";
import nodePath from "node:path";
import { pathToFileURL } from "node:url";
import Ajv2020 from "ajv/dist/2020.js";
import addFormats from "ajv-formats";
import { glob } from "glob";

const ajv = new Ajv2020({
  strict: true,
  allErrors: true,
  validateFormats: true,
});

addFormats(ajv);

const files = (await glob("schema/**/*.schema.json")).sort();

if (files.length === 0) {
  throw new Error("No schema files found under schema/");
}

const schemas = await Promise.all(
  files.map(async (schemaPath) => {
    const text = await readFile(schemaPath, "utf8");
    return {
      path: schemaPath,
      fileUrl: pathToFileURL(nodePath.resolve(schemaPath)).href,
      schema: JSON.parse(text),
    };
  }),
);

for (const { path, fileUrl, schema } of schemas) {
  console.info(`Loading ${path}`);
  ajv.addSchema(schema, fileUrl);
}

for (const { path, fileUrl, schema } of schemas) {
  console.info(`Validating ${path}`);
  ajv.getSchema(fileUrl);
}

console.info(`Validated ${schemas.length} schema files with strict draft 2020-12.`);
