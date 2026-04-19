import path from 'node:path'
import { fileURLToPath } from 'node:url'

import js from '@eslint/js'
import tsPlugin from '@typescript-eslint/eslint-plugin'
import tsParser from '@typescript-eslint/parser'
import nextCoreWebVitals from 'eslint-config-next/core-web-vitals'
import prettierConfig from 'eslint-config-prettier'
import importPlugin from 'eslint-plugin-import'
import prettierPlugin from 'eslint-plugin-prettier'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

const sharedRestrictedSyntax = {
  unknownKeyword: {
    selector: 'TSUnknownKeyword',
    message:
      "The 'unknown' type is forbidden here. Use cast utilities from '@/utils/cast' instead.",
  },
  typeAlias: {
    selector: 'TSTypeAliasDeclaration',
    message: 'Type definitions must be in types/ folder',
  },
  interface: {
    selector: 'TSInterfaceDeclaration',
    message: 'Interface definitions must be in interfaces/ folder',
  },
  exportAll: {
    selector: 'ExportAllDeclaration',
    message:
      'Re-export statements are forbidden. Import and re-export explicitly or use direct imports.',
  },
  exportNamed: {
    selector: 'ExportNamedDeclaration[source]',
    message:
      'Re-export statements are forbidden. Import and re-export explicitly or use direct imports.',
  },
}

export default [
  {
    ignores: [
      'dist/',
      'node_modules/',
      '*.js',
      '**/*.d.ts',
      '.next/',
      'eslint.config.mjs',
    ],
  },
  js.configs.recommended,
  ...nextCoreWebVitals,
  {
    languageOptions: {
      parser: tsParser,
      parserOptions: {
        ecmaVersion: 2022,
        sourceType: 'module',
        project: './tsconfig.json',
        tsconfigRootDir: __dirname,
      },
      globals: {
        JSX: 'readonly',
        React: 'readonly',
      },
    },
    plugins: {
      '@typescript-eslint': tsPlugin,
      prettier: prettierPlugin,
      import: importPlugin,
    },
    rules: {
      ...tsPlugin.configs.recommended.rules,
      ...tsPlugin.configs['recommended-requiring-type-checking'].rules,
      ...prettierConfig.rules,
      'no-warning-comments': [
        'error',
        { terms: ['eslint-disable'], location: 'start' },
      ],
      'prettier/prettier': 'error',
      'no-unused-vars': 'off',
      '@typescript-eslint/no-unused-vars': [
        'error',
        {
          argsIgnorePattern: '^_',
          varsIgnorePattern: '^_',
          caughtErrorsIgnorePattern: '^_',
        },
      ],
      '@typescript-eslint/no-explicit-any': 'error',
      '@typescript-eslint/explicit-function-return-type': 'error',
      '@typescript-eslint/explicit-module-boundary-types': 'error',
      '@typescript-eslint/no-inferrable-types': 'off',
      '@typescript-eslint/no-non-null-assertion': 'error',
      '@typescript-eslint/naming-convention': [
        'error',
        {
          selector: 'interface',
          format: ['PascalCase'],
          prefix: ['I'],
          filter: {
            regex: '^(Assertion|AsymmetricMatchersContaining)$',
            match: false,
          },
        },
      ],
      complexity: ['error', { max: 10 }],
      'max-classes-per-file': ['error', 1],
      'max-depth': ['error', 4],
      'max-len': [
        'error',
        { code: 80, tabWidth: 2, ignoreUrls: true, ignoreStrings: true },
      ],
      'max-lines': [
        'error',
        { max: 300, skipBlankLines: true, skipComments: true },
      ],
      'max-lines-per-function': [
        'error',
        { max: 50, skipBlankLines: true, skipComments: true },
      ],
      'max-nested-callbacks': ['error', 3],
      'max-params': ['error', 5],
      'no-console': 'error',
      'no-debugger': 'error',
      'no-alert': 'error',
      'no-var': 'error',
      'prefer-const': 'error',
      'prefer-arrow-callback': 'error',
      'arrow-spacing': 'error',
      'no-duplicate-imports': 'error',
      'import/order': [
        'error',
        {
          groups: [
            'builtin',
            'external',
            'internal',
            'parent',
            'sibling',
            'index',
          ],
          'newlines-between': 'always',
          alphabetize: {
            order: 'asc',
            caseInsensitive: true,
          },
        },
      ],
      'no-restricted-syntax': [
        'error',
        sharedRestrictedSyntax.unknownKeyword,
        sharedRestrictedSyntax.typeAlias,
        sharedRestrictedSyntax.interface,
        sharedRestrictedSyntax.exportAll,
        sharedRestrictedSyntax.exportNamed,
      ],
    },
  },
  {
    files: ['interfaces/**/*.ts', 'interfaces/**/*.tsx'],
    rules: {
      'no-restricted-syntax': [
        'error',
        sharedRestrictedSyntax.unknownKeyword,
        sharedRestrictedSyntax.typeAlias,
        sharedRestrictedSyntax.exportAll,
        sharedRestrictedSyntax.exportNamed,
      ],
    },
  },
  {
    files: ['types/**/*.ts', 'types/**/*.tsx'],
    rules: {
      'no-restricted-syntax': [
        'error',
        sharedRestrictedSyntax.unknownKeyword,
        sharedRestrictedSyntax.interface,
        sharedRestrictedSyntax.exportAll,
        sharedRestrictedSyntax.exportNamed,
      ],
    },
  },
  {
    files: ['app/layout.tsx', 'app/page.tsx'],
    rules: {
      'no-console': 'off',
    },
  },
  {
    files: ['components/**/*.tsx', 'app/**/*.tsx'],
    rules: {
      'max-lines-per-function': [
        'error',
        { max: 100, skipBlankLines: true, skipComments: true },
      ],
    },
  },
]
