import { createBrowserRouter } from 'react-router-dom'
import PluginList from './routes/PluginList'
import PluginDetail from './routes/PluginDetail'
import PluginDiscover from './routes/PluginDiscover'

export const router = createBrowserRouter([
  {
    path: '/',
    element: <PluginList />,
  },
  {
    path: '/plugins/:id',
    element: <PluginDetail />,
  },
  {
    path: '/plugins/discover',
    element: <PluginDiscover />,
  },
])

