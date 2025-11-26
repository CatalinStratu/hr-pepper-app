import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';
import { analyticsApi } from '../services/api';
import type { Analytics } from '../types';

const COLORS = ['#4CAF50', '#2196F3', '#FF9800', '#f44336', '#9C27B0'];

function Dashboard() {
  const [analytics, setAnalytics] = useState<Analytics | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadAnalytics();
  }, []);

  const loadAnalytics = async () => {
    try {
      const data = await analyticsApi.get();
      setAnalytics(data);
    } catch (error) {
      console.error('Failed to load analytics:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="loading">
        <div className="spinner"></div>
      </div>
    );
  }

  const scoreDistributionData = analytics?.score_distribution
    ? Object.entries(analytics.score_distribution).map(([name, value]) => ({ name, value }))
    : [];

  return (
    <div>
      <div className="page-header">
        <h2>Dashboard</h2>
      </div>

      {/* Stats Grid */}
      <div className="stats-grid">
        <div className="stat-card">
          <div className="stat-label">Total Interviews</div>
          <div className="stat-value">{analytics?.total_interviews || 0}</div>
        </div>
        <div className="stat-card">
          <div className="stat-label">Completed</div>
          <div className="stat-value">{analytics?.completed_interviews || 0}</div>
        </div>
        <div className="stat-card">
          <div className="stat-label">Average Score</div>
          <div className="stat-value">
            {analytics?.average_score ? `${(analytics.average_score * 100).toFixed(0)}%` : 'N/A'}
          </div>
        </div>
        <div className="stat-card">
          <div className="stat-label">Avg Duration</div>
          <div className="stat-value">
            {analytics?.average_duration ? `${Math.round(analytics.average_duration / 60)}m` : 'N/A'}
          </div>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: '20px' }}>
        {/* Question Performance Chart */}
        <div className="card">
          <div className="card-header">
            <h3 className="card-title">Question Performance</h3>
          </div>
          <div className="chart-container">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={analytics?.question_stats || []}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="question_id" />
                <YAxis domain={[0, 1]} tickFormatter={(v) => `${(v * 100).toFixed(0)}%`} />
                <Tooltip formatter={(value: number) => `${(value * 100).toFixed(0)}%`} />
                <Bar dataKey="average_score" fill="#1976D2" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Score Distribution */}
        <div className="card">
          <div className="card-header">
            <h3 className="card-title">Score Distribution</h3>
          </div>
          <div className="chart-container">
            {scoreDistributionData.length > 0 ? (
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={scoreDistributionData}
                    cx="50%"
                    cy="50%"
                    innerRadius={60}
                    outerRadius={80}
                    paddingAngle={5}
                    dataKey="value"
                    label={({ name, value }) => `${name}: ${value}`}
                  >
                    {scoreDistributionData.map((_, index) => (
                      <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip />
                </PieChart>
              </ResponsiveContainer>
            ) : (
              <div className="empty-state">
                <p>No data available yet</p>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Recent Interviews */}
      <div className="card">
        <div className="card-header">
          <h3 className="card-title">Recent Interviews</h3>
          <Link to="/interviews" className="btn btn-secondary">View All</Link>
        </div>
        {analytics?.recent_interviews && analytics.recent_interviews.length > 0 ? (
          <table className="table">
            <thead>
              <tr>
                <th>Candidate</th>
                <th>Position</th>
                <th>Status</th>
                <th>Score</th>
                <th>Action</th>
              </tr>
            </thead>
            <tbody>
              {analytics.recent_interviews.slice(0, 5).map((interview) => (
                <tr key={interview.id}>
                  <td>{interview.candidate_name}</td>
                  <td>{interview.position}</td>
                  <td>
                    <span className={`status-badge ${interview.status}`}>
                      {interview.status.replace('_', ' ')}
                    </span>
                  </td>
                  <td>
                    {interview.total_score !== undefined && interview.total_score !== null
                      ? `${(interview.total_score * 100).toFixed(0)}%`
                      : '-'}
                  </td>
                  <td>
                    <Link to={`/interviews/${interview.id}`} className="btn btn-secondary" style={{ padding: '6px 12px' }}>
                      View
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <div className="empty-state">
            <h3>No interviews yet</h3>
            <p>Interviews conducted by Pepper will appear here</p>
          </div>
        )}
      </div>
    </div>
  );
}

export default Dashboard;
