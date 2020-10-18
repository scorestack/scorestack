module.exports = {
  root: true,
  extends: ['@elastic/eslint-config-kibana', 'plugin:@elastic/eui/recommended'],
  rules: {
    '@kbn/eslint/require-license-header': 'off',
    '@typescript-eslint/camelcase': 'off',
    '@typescript-eslint/class-name-casing': 'off',
  },
};
