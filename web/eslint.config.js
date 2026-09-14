// ESLint 9 flat config：TypeScript + React 规则基线。
import js from '@eslint/js'
import tseslint from 'typescript-eslint'
import react from 'eslint-plugin-react'
import reactHooks from 'eslint-plugin-react-hooks'

export default tseslint.config(
  { ignores: ['dist', 'node_modules', '../internal/webui/dist'] },
  js.configs.recommended,
  ...tseslint.configs.recommended,
  {
    files: ['src/**/*.{ts,tsx}'],
    plugins: { react, 'react-hooks': reactHooks },
    languageOptions: {
      parserOptions: { ecmaFeatures: { jsx: true } },
    },
    settings: { react: { version: 'detect' } },
    rules: {
      ...react.configs.recommended.rules,
      ...reactHooks.configs.recommended.rules,
      // 本项目使用 react-jsx 转换，无需在作用域内引入 React。
      'react/react-in-jsx-scope': 'off',
      'react/prop-types': 'off',
      // 数据获取模式（effect 内调用 async 函数再 setState）是本项目页面的标准写法，
      // react-hooks v7 的 set-state-in-effect 对此误报，关闭。
      'react-hooks/set-state-in-effect': 'off',
    },
  },
)
