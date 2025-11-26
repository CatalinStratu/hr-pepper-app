import { Routes, Route, NavLink } from 'react-router-dom';
import Dashboard from './pages/Dashboard';
import Interviews from './pages/Interviews';
import InterviewDetail from './pages/InterviewDetail';
import Questions from './pages/Questions';

function App() {
  return (
    <div className="app">
      <aside className="sidebar">
        <div className="sidebar-header">
          <h1>HR Interview</h1>
          <p>Pepper Robot System</p>
        </div>
        <ul className="nav-links">
          <li>
            <NavLink to="/" className={({ isActive }) => isActive ? 'active' : ''}>
              <span>📊</span> Dashboard
            </NavLink>
          </li>
          <li>
            <NavLink to="/interviews" className={({ isActive }) => isActive ? 'active' : ''}>
              <span>👥</span> Interviews
            </NavLink>
          </li>
          <li>
            <NavLink to="/questions" className={({ isActive }) => isActive ? 'active' : ''}>
              <span>❓</span> Questions
            </NavLink>
          </li>
        </ul>
      </aside>
      <main className="main-content">
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/interviews" element={<Interviews />} />
          <Route path="/interviews/:id" element={<InterviewDetail />} />
          <Route path="/questions" element={<Questions />} />
        </Routes>
      </main>
    </div>
  );
}

export default App;
