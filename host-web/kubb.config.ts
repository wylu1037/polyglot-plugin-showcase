import { defineConfig } from '@kubb/core'

export default defineConfig({
  input: {
    // 指向后端 Swagger JSON 地址
    path: 'http://localhost:8080/swagger/doc.json',
  },
  output: {
    path: './src/api/generated',
    clean: true,
  },
  plugins: [
    ['@kubb/plugin-oas', { output: false }],
    ['@kubb/plugin-ts', { 
      output: { 
        path: 'types',
      },
    }],
    ['@kubb/plugin-client', { 
      output: { 
        path: 'clients',
      },
      client: {
        importPath: 'axios',
      },
    }],
    ['@kubb/plugin-react-query', { 
      output: { 
        path: 'hooks',
      },
      client: {
        importPath: '../clients/axios',
      },
      query: {
        key: (key) => key,
        methods: ['get'],
      },
      mutation: {
        key: (key) => key,
        methods: ['post', 'put', 'patch', 'delete'],
      },
    }],
  ],
})

