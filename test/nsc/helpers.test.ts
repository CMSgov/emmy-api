import fc from 'fast-check';
import { describe, expect, it } from 'vitest';
import { z } from 'zod';
import { DateFlexibleSchema, SsnSchema } from '../../src/nsc/helpers';

const ssnDigitsArb = fc.array(fc.integer({ min: 1, max: 9 }), {
  minLength: 9,
  maxLength: 9,
});

const ssnDashArb = ssnDigitsArb.map(
  (ds) =>
    ds.slice(0, 3).join('') +
    '-' +
    ds.slice(3, 5).join('') +
    '-' +
    ds.slice(5).join(''),
);

const ssnSpaceArb = ssnDigitsArb.map(
  (ds) =>
    ds.slice(0, 3).join('') +
    ' ' +
    ds.slice(3, 5).join('') +
    ' ' +
    ds.slice(5).join(''),
);

const ssnShortArb = ssnDigitsArb.map((ds) => ds.slice(0, 4).join(''));

const ssnNormalArb = ssnDigitsArb.map((ds) => ds.join(''));

describe('SsnSchema', () => {
  it('parses XXX-XX-XXXX (dashes)', () => {
    // Arrange / Act / Assert
    fc.assert(
      fc.property(ssnDashArb, (s) => {
        const decoded = SsnSchema.decode(s);
        return !decoded.includes('-');
      }),
    );
  });

  it('parses XXX XX XXXX (spaces)', () => {
    fc.assert(
      fc.property(ssnSpaceArb, (s) => {
        const decoded = SsnSchema.decode(s);
        return !decoded.includes('-');
      }),
    );
  });

  it('parses XXXX', () => {
    fc.assert(
      fc.property(ssnShortArb, (s) => {
        const decoded = SsnSchema.decode(s);
        return decoded === s;
      }),
    );
  });

  it('parses XXXXXXXXX', () => {
    fc.assert(
      fc.property(ssnNormalArb, (s) => {
        const decoded = SsnSchema.decode(s);
        return decoded === s;
      }),
    );
  });
});
