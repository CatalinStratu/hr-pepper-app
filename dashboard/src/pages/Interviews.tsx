import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { format } from 'date-fns';
import { interviewsApi } from '../services/api';
import type { Interview } from '../types';

function Interviews() {
  const [interviews, setInterviews] = useState<Interview[]>([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState<string>('');

  useEffect(() => {
    loadInterviews();
  }, [filter]);

  const loadInterviews = async () => {
    try {
      setLoading(true);
      const data = await interviewsApi.getAll(filter || undefined);
      setInterviews(data);
    } catch (error) {
      console.error('Failed to load interviews:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this interview?')) return;

    try {
      await interviewsApi.delete(id);
      loadInterviews();
    } catch (error) {
      console.error('Failed to delete interview:', error);
    }
  };

  const getScoreClass = (score?: number): string => {
    if (score === undefined || score === null) return '';
    if (score >= 0.8) return 'excellent';
    if (score >= 0.6) return 'good';
    if (score >= 0.4) return 'average';
    return 'poor';
  };

  return (
    <div>
      <div className="page-header">
        <h2>Interviews</h2>
        <div style={{ display: 'flex', gap: '12px' }}>
          <select
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            style={{ padding: '10px 14px', borderRadius: '8px', border: '1px solid #ddd' }}
          >
            <option value="">All Status</option>
            <option value="pending">Pending</option>
            <option value="in_progress">In Progress</option>
            <option value="completed">Completed</option>
            <option value="cancelled">Cancelled</option>
          </select>
        </div>
      </div>

      <div className="card">
        {loading ? (
          <div className="loading">
            <div className="spinner"></div>
          </div>
        ) : interviews.length > 0 ? (
          <table className="table">
            <thead>
              <tr>
                <th>Candidate</th>
                <th>Email</th>
                <th>Position</th>
                <th>Status</th>
                <th>Score</th>
                <th>Date</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {interviews.map((interview) => (
                <tr key={interview.id}>
                  <td>
                    <strong>{interview.candidate_name}</strong>
                  </td>
                  <td>{interview.candidate_email || '-'}</td>
                  <td>{interview.position}</td>
                  <td>
                    <span className={`status-badge ${interview.status}`}>
                      {interview.status.replace('_', ' ')}
                    </span>
                  </td>
                  <td>
                    {interview.total_score !== undefined && interview.total_score !== null ? (
                      <span className={`score-badge ${getScoreClass(interview.total_score)}`}>
                        {(interview.total_score * 100).toFixed(0)}%
                      </span>
                    ) : (
                      '-'
                    )}
                  </td>
                  <td>
                    {interview.created_at
                      ? format(new Date(interview.created_at), 'MMM d, yyyy')
                      : '-'}
                  </td>
                  <td>
                    <div style={{ display: 'flex', gap: '8px' }}>
                      <Link
                        to={`/interviews/${interview.id}`}
                        className="btn btn-secondary"
                        style={{ padding: '6px 12px' }}
                      >
                        View
                      </Link>
                      <button
                        onClick={() => handleDelete(interview.id)}
                        className="btn btn-danger"
                        style={{ padding: '6px 12px' }}
                      >
                        Delete
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <div className="empty-state">
            <h3>No interviews found</h3>
            <p>Interviews will appear here once candidates complete them on Pepper</p>
          </div>
        )}
      </div>
    </div>
  );
}

export default Interviews;
