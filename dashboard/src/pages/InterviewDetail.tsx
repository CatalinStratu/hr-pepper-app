import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { format } from 'date-fns';
import { interviewsApi } from '../services/api';
import type { InterviewDetail as InterviewDetailType, Response } from '../types';

function InterviewDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [data, setData] = useState<InterviewDetailType | null>(null);
  const [loading, setLoading] = useState(true);
  const [notes, setNotes] = useState('');
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (id) {
      loadInterview();
    }
  }, [id]);

  const loadInterview = async () => {
    try {
      const result = await interviewsApi.getById(id!);
      setData(result);
      setNotes(result.interview.notes || '');
    } catch (error) {
      console.error('Failed to load interview:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleUpdateNotes = async () => {
    if (!id) return;

    try {
      setSaving(true);
      await interviewsApi.update(id, { notes });
      alert('Notes saved successfully');
    } catch (error) {
      console.error('Failed to save notes:', error);
    } finally {
      setSaving(false);
    }
  };

  const handleUpdateStatus = async (status: string) => {
    if (!id) return;

    try {
      await interviewsApi.update(id, { status });
      loadInterview();
    } catch (error) {
      console.error('Failed to update status:', error);
    }
  };

  const getScoreColor = (score?: number): string => {
    if (score === undefined || score === null) return '#666';
    if (score >= 0.7) return '#4CAF50';
    if (score >= 0.5) return '#FF9800';
    return '#f44336';
  };

  if (loading) {
    return (
      <div className="loading">
        <div className="spinner"></div>
      </div>
    );
  }

  if (!data) {
    return (
      <div className="empty-state">
        <h3>Interview not found</h3>
        <button onClick={() => navigate('/interviews')} className="btn btn-primary">
          Back to Interviews
        </button>
      </div>
    );
  }

  const { interview, responses } = data;

  return (
    <div>
      <div className="page-header">
        <h2>Interview Details</h2>
        <button onClick={() => navigate('/interviews')} className="btn btn-secondary">
          Back to List
        </button>
      </div>

      {/* Candidate Info Card */}
      <div className="card">
        <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: '24px' }}>
          <div>
            <h3 style={{ fontSize: '24px', marginBottom: '8px' }}>{interview.candidate_name}</h3>
            <p style={{ color: '#666', marginBottom: '16px' }}>
              {interview.candidate_email || 'No email provided'}
            </p>
            <div style={{ display: 'flex', gap: '24px', marginBottom: '16px' }}>
              <div>
                <span style={{ color: '#666', fontSize: '13px' }}>Position</span>
                <p style={{ fontWeight: '600' }}>{interview.position}</p>
              </div>
              <div>
                <span style={{ color: '#666', fontSize: '13px' }}>Status</span>
                <p>
                  <span className={`status-badge ${interview.status}`}>
                    {interview.status.replace('_', ' ')}
                  </span>
                </p>
              </div>
              <div>
                <span style={{ color: '#666', fontSize: '13px' }}>Date</span>
                <p style={{ fontWeight: '600' }}>
                  {interview.created_at
                    ? format(new Date(interview.created_at), 'MMMM d, yyyy h:mm a')
                    : '-'}
                </p>
              </div>
            </div>
            <div style={{ display: 'flex', gap: '8px' }}>
              {interview.status !== 'completed' && (
                <button
                  onClick={() => handleUpdateStatus('completed')}
                  className="btn btn-success"
                >
                  Mark Completed
                </button>
              )}
              {interview.status !== 'cancelled' && (
                <button
                  onClick={() => handleUpdateStatus('cancelled')}
                  className="btn btn-danger"
                >
                  Cancel Interview
                </button>
              )}
            </div>
          </div>
          <div style={{ textAlign: 'center', borderLeft: '1px solid #eee', paddingLeft: '24px' }}>
            <span style={{ color: '#666', fontSize: '13px' }}>Overall Score</span>
            <div
              style={{
                fontSize: '64px',
                fontWeight: '700',
                color: getScoreColor(interview.total_score),
                lineHeight: '1.2',
              }}
            >
              {interview.total_score !== undefined && interview.total_score !== null
                ? `${(interview.total_score * 100).toFixed(0)}%`
                : 'N/A'}
            </div>
            <p style={{ color: '#666', fontSize: '13px' }}>
              {responses?.length || 0} responses recorded
            </p>
          </div>
        </div>
      </div>

      {/* Responses */}
      <div className="card">
        <div className="card-header">
          <h3 className="card-title">Interview Responses</h3>
        </div>
        {responses && responses.length > 0 ? (
          <div>
            {responses.map((response: Response, index: number) => (
              <div key={response.id} className="response-card">
                <div className="question">
                  Q{index + 1}: {response.question_text || `Question ${response.question_id}`}
                </div>
                <div className="answer">{response.response_text}</div>
                <div className="metrics">
                  <span>
                    Duration: {Math.floor(response.duration / 60)}m {response.duration % 60}s
                  </span>
                  {response.sentiment_score !== undefined && (
                    <span style={{ color: getScoreColor(response.sentiment_score) }}>
                      Sentiment: {(response.sentiment_score * 100).toFixed(0)}%
                    </span>
                  )}
                  {response.confidence_score !== undefined && (
                    <span style={{ color: getScoreColor(response.confidence_score) }}>
                      Confidence: {(response.confidence_score * 100).toFixed(0)}%
                    </span>
                  )}
                </div>
                {response.analysis && (
                  <div style={{ marginTop: '12px', padding: '12px', background: '#e3f2fd', borderRadius: '8px', fontSize: '13px' }}>
                    <strong>Analysis:</strong> {response.analysis}
                  </div>
                )}
              </div>
            ))}
          </div>
        ) : (
          <div className="empty-state">
            <p>No responses recorded yet</p>
          </div>
        )}
      </div>

      {/* Notes Section */}
      <div className="card">
        <div className="card-header">
          <h3 className="card-title">HR Notes</h3>
        </div>
        <div className="form-group">
          <textarea
            value={notes}
            onChange={(e) => setNotes(e.target.value)}
            rows={4}
            placeholder="Add notes about this candidate..."
            style={{ resize: 'vertical' }}
          />
        </div>
        <button onClick={handleUpdateNotes} className="btn btn-primary" disabled={saving}>
          {saving ? 'Saving...' : 'Save Notes'}
        </button>
      </div>
    </div>
  );
}

export default InterviewDetail;
