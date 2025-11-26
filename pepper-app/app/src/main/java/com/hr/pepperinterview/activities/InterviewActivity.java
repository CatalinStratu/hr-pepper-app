package com.hr.pepperinterview.activities;

import android.content.Intent;
import android.os.Bundle;
import android.os.CountDownTimer;
import android.os.Handler;
import android.os.Looper;
import android.speech.RecognitionListener;
import android.speech.RecognizerIntent;
import android.speech.SpeechRecognizer;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.ProgressBar;
import android.widget.TextView;
import android.widget.Toast;

import androidx.appcompat.app.AppCompatActivity;
import androidx.cardview.widget.CardView;

import com.aldebaran.qi.sdk.QiContext;
import com.aldebaran.qi.sdk.QiSDK;
import com.aldebaran.qi.sdk.RobotLifecycleCallbacks;
import com.aldebaran.qi.sdk.builder.ListenBuilder;
import com.aldebaran.qi.sdk.builder.PhraseSetBuilder;
import com.aldebaran.qi.sdk.builder.SayBuilder;
import com.aldebaran.qi.sdk.object.conversation.Listen;
import com.aldebaran.qi.sdk.object.conversation.ListenResult;
import com.aldebaran.qi.sdk.object.conversation.PhraseSet;
import com.aldebaran.qi.sdk.object.conversation.Say;
import com.hr.pepperinterview.R;
import com.hr.pepperinterview.api.ApiClient;
import com.hr.pepperinterview.models.Interview;
import com.hr.pepperinterview.models.Question;
import com.hr.pepperinterview.models.Response;

import java.util.ArrayList;
import java.util.List;
import java.util.Locale;

import retrofit2.Call;
import retrofit2.Callback;

public class InterviewActivity extends AppCompatActivity implements RobotLifecycleCallbacks {
    private static final String TAG = "InterviewActivity";

    private TextView tvQuestionNumber;
    private TextView tvQuestion;
    private TextView tvCategory;
    private TextView tvTimer;
    private EditText editResponse;
    private Button btnRecord;
    private Button btnSubmit;
    private Button btnSkip;
    private ProgressBar progressBar;
    private CardView cardFeedback;
    private TextView tvFeedback;

    private String interviewId;
    private String candidateName;
    private String position;
    private List<Question> questions = new ArrayList<>();
    private int currentQuestionIndex = 0;
    private long questionStartTime;

    private QiContext qiContext;
    private SpeechRecognizer speechRecognizer;
    private boolean isListening = false;
    private CountDownTimer countDownTimer;

    private double totalSentimentScore = 0;
    private int responsesSubmitted = 0;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_interview);

        // Get intent data
        interviewId = getIntent().getStringExtra("interview_id");
        candidateName = getIntent().getStringExtra("candidate_name");
        position = getIntent().getStringExtra("position");

        // Initialize views
        tvQuestionNumber = findViewById(R.id.tvQuestionNumber);
        tvQuestion = findViewById(R.id.tvQuestion);
        tvCategory = findViewById(R.id.tvCategory);
        tvTimer = findViewById(R.id.tvTimer);
        editResponse = findViewById(R.id.editResponse);
        btnRecord = findViewById(R.id.btnRecord);
        btnSubmit = findViewById(R.id.btnSubmit);
        btnSkip = findViewById(R.id.btnSkip);
        progressBar = findViewById(R.id.progressBar);
        cardFeedback = findViewById(R.id.cardFeedback);
        tvFeedback = findViewById(R.id.tvFeedback);

        // Setup speech recognizer
        setupSpeechRecognizer();

        // Register for robot lifecycle
        QiSDK.register(this, this);

        // Set click listeners
        btnRecord.setOnClickListener(v -> toggleRecording());
        btnSubmit.setOnClickListener(v -> submitResponse());
        btnSkip.setOnClickListener(v -> skipQuestion());

        // Load questions from server
        loadQuestions();
    }

    @Override
    protected void onDestroy() {
        QiSDK.unregister(this, this);
        if (speechRecognizer != null) {
            speechRecognizer.destroy();
        }
        if (countDownTimer != null) {
            countDownTimer.cancel();
        }
        super.onDestroy();
    }

    @Override
    public void onRobotFocusGained(QiContext qiContext) {
        this.qiContext = qiContext;
        Log.i(TAG, "Robot focus gained");
    }

    @Override
    public void onRobotFocusLost() {
        this.qiContext = null;
        Log.i(TAG, "Robot focus lost");
    }

    @Override
    public void onRobotFocusRefused(String reason) {
        Log.e(TAG, "Robot focus refused: " + reason);
    }

    private void setupSpeechRecognizer() {
        if (SpeechRecognizer.isRecognitionAvailable(this)) {
            speechRecognizer = SpeechRecognizer.createSpeechRecognizer(this);
            speechRecognizer.setRecognitionListener(new RecognitionListener() {
                @Override
                public void onReadyForSpeech(Bundle params) {
                    Log.d(TAG, "Ready for speech");
                }

                @Override
                public void onBeginningOfSpeech() {
                    Log.d(TAG, "Beginning of speech");
                }

                @Override
                public void onRmsChanged(float rmsdB) {}

                @Override
                public void onBufferReceived(byte[] buffer) {}

                @Override
                public void onEndOfSpeech() {
                    Log.d(TAG, "End of speech");
                    runOnUiThread(() -> {
                        isListening = false;
                        btnRecord.setText("🎤 Record Answer");
                    });
                }

                @Override
                public void onError(int error) {
                    Log.e(TAG, "Speech recognition error: " + error);
                    runOnUiThread(() -> {
                        isListening = false;
                        btnRecord.setText("🎤 Record Answer");
                    });
                }

                @Override
                public void onResults(Bundle results) {
                    ArrayList<String> matches = results.getStringArrayList(SpeechRecognizer.RESULTS_RECOGNITION);
                    if (matches != null && !matches.isEmpty()) {
                        String text = editResponse.getText().toString();
                        if (!text.isEmpty()) {
                            text += " ";
                        }
                        text += matches.get(0);
                        editResponse.setText(text);
                    }
                }

                @Override
                public void onPartialResults(Bundle partialResults) {}

                @Override
                public void onEvent(int eventType, Bundle params) {}
            });
        }
    }

    private void loadQuestions() {
        progressBar.setVisibility(View.VISIBLE);

        // The questions were already returned when creating the interview
        // Fetch interview details to get questions
        ApiClient.getInstance().getApiService().getInterview(interviewId)
                .enqueue(new Callback<Interview>() {
                    @Override
                    public void onResponse(Call<Interview> call, retrofit2.Response<Interview> response) {
                        // For now, use hardcoded questions as a fallback
                        // In production, these come from the server response
                        questions = getDefaultQuestions();
                        progressBar.setVisibility(View.GONE);
                        displayCurrentQuestion();
                    }

                    @Override
                    public void onFailure(Call<Interview> call, Throwable t) {
                        // Use default questions
                        questions = getDefaultQuestions();
                        progressBar.setVisibility(View.GONE);
                        displayCurrentQuestion();
                    }
                });
    }

    private List<Question> getDefaultQuestions() {
        List<Question> defaultQuestions = new ArrayList<>();

        Question q1 = new Question();
        q1.setId("q1");
        q1.setText("Tell me about yourself and your background.");
        q1.setCategory("introduction");
        q1.setTimeLimit(120);
        defaultQuestions.add(q1);

        Question q2 = new Question();
        q2.setId("q2");
        q2.setText("What interests you about this position?");
        q2.setCategory("motivation");
        q2.setTimeLimit(90);
        defaultQuestions.add(q2);

        Question q3 = new Question();
        q3.setId("q3");
        q3.setText("Describe a challenging situation you faced at work and how you handled it.");
        q3.setCategory("behavioral");
        q3.setTimeLimit(180);
        defaultQuestions.add(q3);

        Question q4 = new Question();
        q4.setId("q4");
        q4.setText("What are your greatest strengths?");
        q4.setCategory("self-assessment");
        q4.setTimeLimit(90);
        defaultQuestions.add(q4);

        Question q5 = new Question();
        q5.setId("q5");
        q5.setText("Where do you see yourself in 5 years?");
        q5.setCategory("goals");
        q5.setTimeLimit(120);
        defaultQuestions.add(q5);

        return defaultQuestions;
    }

    private void displayCurrentQuestion() {
        if (currentQuestionIndex >= questions.size()) {
            finishInterview();
            return;
        }

        Question question = questions.get(currentQuestionIndex);

        tvQuestionNumber.setText("Question " + (currentQuestionIndex + 1) + " of " + questions.size());
        tvQuestion.setText(question.getText());
        tvCategory.setText("Category: " + question.getCategory());
        editResponse.setText("");
        cardFeedback.setVisibility(View.GONE);

        questionStartTime = System.currentTimeMillis();

        // Start timer
        startTimer(question.getTimeLimit());

        // Robot speaks the question
        speakQuestion(question.getText());
    }

    private void speakQuestion(String questionText) {
        if (qiContext != null) {
            new Thread(() -> {
                try {
                    Say say = SayBuilder.with(qiContext)
                            .withText(questionText)
                            .build();
                    say.run();
                } catch (Exception e) {
                    Log.e(TAG, "Error speaking question: " + e.getMessage());
                }
            }).start();
        }
    }

    private void startTimer(int seconds) {
        if (countDownTimer != null) {
            countDownTimer.cancel();
        }

        countDownTimer = new CountDownTimer(seconds * 1000L, 1000) {
            @Override
            public void onTick(long millisUntilFinished) {
                int secondsRemaining = (int) (millisUntilFinished / 1000);
                int minutes = secondsRemaining / 60;
                int secs = secondsRemaining % 60;
                tvTimer.setText(String.format(Locale.getDefault(), "%02d:%02d", minutes, secs));

                // Change color when time is running low
                if (secondsRemaining <= 30) {
                    tvTimer.setTextColor(getResources().getColor(android.R.color.holo_red_dark));
                } else {
                    tvTimer.setTextColor(getResources().getColor(android.R.color.black));
                }
            }

            @Override
            public void onFinish() {
                tvTimer.setText("00:00");
                Toast.makeText(InterviewActivity.this, "Time's up!", Toast.LENGTH_SHORT).show();

                // Auto-submit if there's a response
                if (!editResponse.getText().toString().trim().isEmpty()) {
                    submitResponse();
                }
            }
        };
        countDownTimer.start();
    }

    private void toggleRecording() {
        if (isListening) {
            stopRecording();
        } else {
            startRecording();
        }
    }

    private void startRecording() {
        if (speechRecognizer != null) {
            Intent intent = new Intent(RecognizerIntent.ACTION_RECOGNIZE_SPEECH);
            intent.putExtra(RecognizerIntent.EXTRA_LANGUAGE_MODEL, RecognizerIntent.LANGUAGE_MODEL_FREE_FORM);
            intent.putExtra(RecognizerIntent.EXTRA_LANGUAGE, Locale.getDefault());
            intent.putExtra(RecognizerIntent.EXTRA_PARTIAL_RESULTS, true);

            speechRecognizer.startListening(intent);
            isListening = true;
            btnRecord.setText("⏹ Stop Recording");
        } else {
            Toast.makeText(this, "Speech recognition not available", Toast.LENGTH_SHORT).show();
        }
    }

    private void stopRecording() {
        if (speechRecognizer != null) {
            speechRecognizer.stopListening();
            isListening = false;
            btnRecord.setText("🎤 Record Answer");
        }
    }

    private void submitResponse() {
        String responseText = editResponse.getText().toString().trim();

        if (responseText.isEmpty()) {
            Toast.makeText(this, "Please provide an answer", Toast.LENGTH_SHORT).show();
            return;
        }

        if (countDownTimer != null) {
            countDownTimer.cancel();
        }

        Question question = questions.get(currentQuestionIndex);
        int duration = (int) ((System.currentTimeMillis() - questionStartTime) / 1000);

        progressBar.setVisibility(View.VISIBLE);
        btnSubmit.setEnabled(false);

        Response.SubmitResponseRequest request = new Response.SubmitResponseRequest(
                interviewId, question.getId(), responseText, duration
        );

        ApiClient.getInstance().getApiService().submitResponse(request)
                .enqueue(new Callback<Response.SubmitResponseResult>() {
                    @Override
                    public void onResponse(Call<Response.SubmitResponseResult> call,
                                           retrofit2.Response<Response.SubmitResponseResult> response) {
                        progressBar.setVisibility(View.GONE);
                        btnSubmit.setEnabled(true);

                        if (response.isSuccessful() && response.body() != null) {
                            Response.SubmitResponseResult result = response.body();
                            Response.AnalysisResult analysis = result.getAnalysis();

                            if (analysis != null) {
                                totalSentimentScore += analysis.getSentimentScore();
                                responsesSubmitted++;

                                // Show feedback
                                cardFeedback.setVisibility(View.VISIBLE);
                                tvFeedback.setText(analysis.getFeedback());

                                // Robot gives brief feedback
                                if (qiContext != null) {
                                    new Thread(() -> {
                                        try {
                                            String feedback = analysis.getSentimentScore() >= 0.6
                                                    ? "Good answer! Let's move to the next question."
                                                    : "Thank you for your response. Let's continue.";
                                            Say say = SayBuilder.with(qiContext)
                                                    .withText(feedback)
                                                    .build();
                                            say.run();
                                        } catch (Exception e) {
                                            Log.e(TAG, "Error speaking feedback: " + e.getMessage());
                                        }
                                    }).start();
                                }
                            }

                            // Move to next question after delay
                            new Handler(Looper.getMainLooper()).postDelayed(() -> {
                                currentQuestionIndex++;
                                displayCurrentQuestion();
                            }, 3000);

                        } else {
                            Toast.makeText(InterviewActivity.this,
                                    "Failed to submit response", Toast.LENGTH_SHORT).show();
                        }
                    }

                    @Override
                    public void onFailure(Call<Response.SubmitResponseResult> call, Throwable t) {
                        progressBar.setVisibility(View.GONE);
                        btnSubmit.setEnabled(true);
                        Log.e(TAG, "Submit response failed: " + t.getMessage());
                        Toast.makeText(InterviewActivity.this,
                                "Connection error: " + t.getMessage(), Toast.LENGTH_SHORT).show();
                    }
                });
    }

    private void skipQuestion() {
        if (countDownTimer != null) {
            countDownTimer.cancel();
        }
        currentQuestionIndex++;
        displayCurrentQuestion();
    }

    private void finishInterview() {
        // Calculate final score
        double finalScore = responsesSubmitted > 0 ? totalSentimentScore / responsesSubmitted : 0;

        // Update interview status on server
        Interview.UpdateInterviewRequest updateRequest = new Interview.UpdateInterviewRequest(
                "completed", finalScore
        );

        ApiClient.getInstance().getApiService().updateInterview(interviewId, updateRequest)
                .enqueue(new Callback<Void>() {
                    @Override
                    public void onResponse(Call<Void> call, retrofit2.Response<Void> response) {
                        Log.d(TAG, "Interview completed successfully");
                    }

                    @Override
                    public void onFailure(Call<Void> call, Throwable t) {
                        Log.e(TAG, "Failed to update interview status: " + t.getMessage());
                    }
                });

        // Robot says goodbye
        if (qiContext != null) {
            new Thread(() -> {
                try {
                    Say say = SayBuilder.with(qiContext)
                            .withText("Thank you " + candidateName + " for completing the interview! We will review your responses and get back to you soon. Good luck!")
                            .build();
                    say.run();
                } catch (Exception e) {
                    Log.e(TAG, "Error with goodbye speech: " + e.getMessage());
                }
            }).start();
        }

        // Go to completion screen
        Intent intent = new Intent(this, CompleteActivity.class);
        intent.putExtra("candidate_name", candidateName);
        intent.putExtra("questions_answered", responsesSubmitted);
        intent.putExtra("total_questions", questions.size());
        intent.putExtra("score", finalScore);
        startActivity(intent);
        finish();
    }
}
