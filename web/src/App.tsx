import { Dashboard } from './components/Dashboard';

function App() {
  // Public Dashboard: No Auth Required
  return <Dashboard onLogout={() => { }} />;
}

export default App;
