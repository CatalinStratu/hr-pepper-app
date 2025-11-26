import { useState, useEffect } from 'react';
import { questionsApi } from '../services/api';
import type { Question } from '../types';

function Questions() {
  const [questions, setQuestions] = useState<Question[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editingQuestion, setEditingQuestion] = useState<Question | null>(null);
  const [formData, setFormData] = useState({
    text: '',
    category: '',
    difficulty: 'medium',
    time_limit: 120,
    order: 0,
  });

  useEffect(() => {
    loadQuestions();
  }, []);

  const loadQuestions = async () => {
    try {
      const data = await questionsApi.getAll();
      setQuestions(data);
    } catch (error) {
      console.error('Failed to load questions:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleOpenModal = (question?: Question) => {
    if (question) {
      setEditingQuestion(question);
      setFormData({
        text: question.text,
        category: question.category,
        difficulty: question.difficulty,
        time_limit: question.time_limit,
        order: question.order,
      });
    } else {
      setEditingQuestion(null);
      setFormData({
        text: '',
        category: '',
        difficulty: 'medium',
        time_limit: 120,
        order: questions.length + 1,
      });
    }
    setShowModal(true);
  };

  const handleCloseModal = () => {
    setShowModal(false);
    setEditingQuestion(null);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      if (editingQuestion) {
        await questionsApi.update(editingQuestion.id, formData);
      } else {
        await questionsApi.create(formData);
      }
      loadQuestions();
      handleCloseModal();
    } catch (error) {
      console.error('Failed to save question:', error);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this question?')) return;

    try {
      await questionsApi.delete(id);
      loadQuestions();
    } catch (error) {
      console.error('Failed to delete question:', error);
    }
  };

  const getDifficultyColor = (difficulty: string): string => {
    switch (difficulty) {
      case 'easy': return '#4CAF50';
      case 'medium': return '#FF9800';
      case 'hard': return '#f44336';
      default: return '#666';
    }
  };

  const categories = [
    'introduction',
    'motivation',
    'behavioral',
    'self-assessment',
    'goals',
    'teamwork',
    'stress-management',
    'adaptability',
    'problem-solving',
    'technical',
    'closing',
  ];

  return (
    <div>
      <div className="page-header">
        <h2>Interview Questions</h2>
        <button onClick={() => handleOpenModal()} className="btn btn-primary">
          Add Question
        </button>
      </div>

      <div className="card">
        {loading ? (
          <div className="loading">
            <div className="spinner"></div>
          </div>
        ) : questions.length > 0 ? (
          <table className="table">
            <thead>
              <tr>
                <th style={{ width: '50px' }}>#</th>
                <th>Question</th>
                <th>Category</th>
                <th>Difficulty</th>
                <th>Time Limit</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {questions
                .sort((a, b) => a.order - b.order)
                .map((question) => (
                  <tr key={question.id}>
                    <td style={{ fontWeight: '600', color: '#666' }}>{question.order}</td>
                    <td style={{ maxWidth: '400px' }}>{question.text}</td>
                    <td>
                      <span
                        style={{
                          background: '#e3f2fd',
                          padding: '4px 12px',
                          borderRadius: '16px',
                          fontSize: '12px',
                        }}
                      >
                        {question.category}
                      </span>
                    </td>
                    <td>
                      <span
                        style={{
                          color: getDifficultyColor(question.difficulty),
                          fontWeight: '500',
                          textTransform: 'capitalize',
                        }}
                      >
                        {question.difficulty}
                      </span>
                    </td>
                    <td>{question.time_limit}s</td>
                    <td>
                      <div style={{ display: 'flex', gap: '8px' }}>
                        <button
                          onClick={() => handleOpenModal(question)}
                          className="btn btn-secondary"
                          style={{ padding: '6px 12px' }}
                        >
                          Edit
                        </button>
                        <button
                          onClick={() => handleDelete(question.id)}
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
            <h3>No questions found</h3>
            <p>Add interview questions for Pepper to ask candidates</p>
          </div>
        )}
      </div>

      {/* Modal */}
      {showModal && (
        <div className="modal-overlay" onClick={handleCloseModal}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>{editingQuestion ? 'Edit Question' : 'Add New Question'}</h3>
              <button className="modal-close" onClick={handleCloseModal}>
                &times;
              </button>
            </div>
            <form onSubmit={handleSubmit}>
              <div className="form-group">
                <label>Question Text</label>
                <textarea
                  value={formData.text}
                  onChange={(e) => setFormData({ ...formData, text: e.target.value })}
                  rows={3}
                  required
                  placeholder="Enter the interview question..."
                />
              </div>
              <div className="form-group">
                <label>Category</label>
                <select
                  value={formData.category}
                  onChange={(e) => setFormData({ ...formData, category: e.target.value })}
                  required
                >
                  <option value="">Select category</option>
                  {categories.map((cat) => (
                    <option key={cat} value={cat}>
                      {cat.charAt(0).toUpperCase() + cat.slice(1).replace('-', ' ')}
                    </option>
                  ))}
                </select>
              </div>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '16px' }}>
                <div className="form-group">
                  <label>Difficulty</label>
                  <select
                    value={formData.difficulty}
                    onChange={(e) => setFormData({ ...formData, difficulty: e.target.value })}
                  >
                    <option value="easy">Easy</option>
                    <option value="medium">Medium</option>
                    <option value="hard">Hard</option>
                  </select>
                </div>
                <div className="form-group">
                  <label>Time Limit (seconds)</label>
                  <input
                    type="number"
                    value={formData.time_limit}
                    onChange={(e) => setFormData({ ...formData, time_limit: parseInt(e.target.value) })}
                    min={30}
                    max={600}
                  />
                </div>
                <div className="form-group">
                  <label>Order</label>
                  <input
                    type="number"
                    value={formData.order}
                    onChange={(e) => setFormData({ ...formData, order: parseInt(e.target.value) })}
                    min={1}
                  />
                </div>
              </div>
              <div className="modal-actions">
                <button type="button" onClick={handleCloseModal} className="btn btn-secondary">
                  Cancel
                </button>
                <button type="submit" className="btn btn-primary">
                  {editingQuestion ? 'Update Question' : 'Add Question'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

export default Questions;
