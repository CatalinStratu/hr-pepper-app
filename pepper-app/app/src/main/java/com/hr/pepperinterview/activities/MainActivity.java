package com.hr.pepperinterview.activities;

import android.content.Intent;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.ProgressBar;
import android.widget.Toast;

import androidx.appcompat.app.AppCompatActivity;

import com.aldebaran.qi.sdk.QiContext;
import com.aldebaran.qi.sdk.QiSDK;
import com.aldebaran.qi.sdk.RobotLifecycleCallbacks;
import com.aldebaran.qi.sdk.builder.AnimateBuilder;
import com.aldebaran.qi.sdk.builder.AnimationBuilder;
import com.aldebaran.qi.sdk.builder.SayBuilder;
import com.aldebaran.qi.sdk.object.actuation.Animate;
import com.aldebaran.qi.sdk.object.actuation.Animation;
import com.aldebaran.qi.sdk.object.conversation.Say;
import com.hr.pepperinterview.R;
import com.hr.pepperinterview.api.ApiClient;
import com.hr.pepperinterview.models.Interview;
import com.hr.pepperinterview.models.Question;

import java.util.ArrayList;
import java.util.List;

import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;

public class MainActivity extends AppCompatActivity implements RobotLifecycleCallbacks {
    private static final String TAG = "MainActivity";

    private EditText editCandidateName;
    private EditText editCandidateEmail;
    private EditText editPosition;
    private EditText editServerUrl;
    private Button btnStartInterview;
    private ProgressBar progressBar;

    private QiContext qiContext;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        // Initialize views
        editCandidateName = findViewById(R.id.editCandidateName);
        editCandidateEmail = findViewById(R.id.editCandidateEmail);
        editPosition = findViewById(R.id.editPosition);
        editServerUrl = findViewById(R.id.editServerUrl);
        btnStartInterview = findViewById(R.id.btnStartInterview);
        progressBar = findViewById(R.id.progressBar);

        // Set default server URL
        editServerUrl.setText(ApiClient.getInstance().getServerUrl());

        // Register for robot lifecycle
        QiSDK.register(this, this);

        btnStartInterview.setOnClickListener(v -> startInterview());
    }

    @Override
    protected void onDestroy() {
        QiSDK.unregister(this, this);
        super.onDestroy();
    }

    @Override
    public void onRobotFocusGained(QiContext qiContext) {
        this.qiContext = qiContext;
        Log.i(TAG, "Robot focus gained");

        // Greet the candidate
        runOnUiThread(() -> {
            try {
                Say say = SayBuilder.with(qiContext)
                        .withText("Hello! Welcome to the HR Interview. Please enter your details on the screen to begin.")
                        .build();
                say.run();
            } catch (Exception e) {
                Log.e(TAG, "Error greeting: " + e.getMessage());
            }
        });
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

    private void startInterview() {
        String candidateName = editCandidateName.getText().toString().trim();
        String candidateEmail = editCandidateEmail.getText().toString().trim();
        String position = editPosition.getText().toString().trim();
        String serverUrl = editServerUrl.getText().toString().trim();

        // Validate input
        if (candidateName.isEmpty()) {
            editCandidateName.setError("Please enter your name");
            return;
        }
        if (position.isEmpty()) {
            editPosition.setError("Please enter the position");
            return;
        }

        // Update server URL if changed
        if (!serverUrl.equals(ApiClient.getInstance().getServerUrl())) {
            ApiClient.initialize(serverUrl);
        }

        // Show loading
        progressBar.setVisibility(View.VISIBLE);
        btnStartInterview.setEnabled(false);

        // Create interview on server
        Interview.CreateInterviewRequest request = new Interview.CreateInterviewRequest(
                candidateName, candidateEmail, position
        );

        ApiClient.getInstance().getApiService().createInterview(request)
                .enqueue(new Callback<Interview.CreateInterviewResponse>() {
                    @Override
                    public void onResponse(Call<Interview.CreateInterviewResponse> call,
                                           Response<Interview.CreateInterviewResponse> response) {
                        progressBar.setVisibility(View.GONE);
                        btnStartInterview.setEnabled(true);

                        if (response.isSuccessful() && response.body() != null) {
                            Interview.CreateInterviewResponse result = response.body();
                            Interview interview = result.getInterview();
                            List<Question> questions = result.getQuestions();

                            // Robot announces start
                            if (qiContext != null) {
                                try {
                                    Say say = SayBuilder.with(qiContext)
                                            .withText("Great " + candidateName + "! Let's begin your interview for the " + position + " position. Good luck!")
                                            .build();
                                    say.run();
                                } catch (Exception e) {
                                    Log.e(TAG, "Error with robot speech: " + e.getMessage());
                                }
                            }

                            // Start interview activity
                            Intent intent = new Intent(MainActivity.this, InterviewActivity.class);
                            intent.putExtra("interview_id", interview.getId());
                            intent.putExtra("candidate_name", candidateName);
                            intent.putExtra("position", position);
                            intent.putParcelableArrayListExtra("questions", new ArrayList<>()); // We'll fetch in next activity
                            startActivity(intent);

                        } else {
                            Toast.makeText(MainActivity.this,
                                    "Failed to create interview. Please check server connection.",
                                    Toast.LENGTH_LONG).show();
                        }
                    }

                    @Override
                    public void onFailure(Call<Interview.CreateInterviewResponse> call, Throwable t) {
                        progressBar.setVisibility(View.GONE);
                        btnStartInterview.setEnabled(true);
                        Log.e(TAG, "API call failed: " + t.getMessage());
                        Toast.makeText(MainActivity.this,
                                "Connection failed: " + t.getMessage(),
                                Toast.LENGTH_LONG).show();
                    }
                });
    }
}
