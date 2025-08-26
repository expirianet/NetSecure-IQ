module.exports = {
  root: true,
  env: { browser: true, node: true, es6: true },
  parser: 'vue-eslint-parser',
  parserOptions: {
    ecmaVersion: 2021,
    sourceType: 'module',
    extraFileExtensions: ['.vue'],
    parser: '@babel/eslint-parser',
    requireConfigFile: false,
    babelOptions: { presets: ['@babel/preset-env'] },
  },
  extends: ['eslint:recommended', 'plugin:vue/vue3-essential'],
  rules: {
    'vue/multi-word-component-names': 'off',
    'no-unused-vars': ['warn', { args: 'none', ignoreRestSiblings: true }],
  },
  ignorePatterns: ['dist/', 'public/', 'node_modules/', '../backend/**'],
}
