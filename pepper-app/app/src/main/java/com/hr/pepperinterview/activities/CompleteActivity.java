package com.hr.pepperinterview.activities;

import android.content.Intent;
import android.os.Bundle;
import android.util.Log;
import android.widget.Button;
import android.widget.TextView;

import androidx.appcompat.app.AppCompatActivity;

import com.aldebaran.qi.sdk.QiContext;
import com.aldebaran.qi.sdk.QiSDK;
import com.aldebaran.qi.sdk.RobotLifecycleCallbacks;
import com.aldebaran.qi.sdk.builder.SayBuilder;
import com.aldebaran.qi.sdk.object.conversation.Say;
import com.hr.pepperinterview.R;

import java.util.Locale;

public class CompleteActivity extends AppCompatActivity implements RobotLifecycleCallbacks {
    private static final String TAG = "CompleteActivity";

    private TextView tvCompletionMessage;
    private TextView tvStats;
    private TextView tvScoreDisplay;
    private Button btnNewInterview;

    private QiContext qiContext;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_complete);

        // Initialize views
        tvCompletionMessage = findViewById(R.id.tvCompletionMessage);
        tvStats = findViewById(R.id.tvStats);
        tvScoreDisplay = findViewById(R.id.tvScoreDisplay);
        btnNewInterview = findViewById(R.id.btnNewInterview);

        // Get intent data
        String candidateName = getIntent().getStringExtra("candidate_name");
        int questionsAnswered = getIntent().getIntExtra("questions_answered", 0);
        int totalQuestions = getIntent().getIntExtra("total_questions", 0);
        double score = getIntent().getDoubleExtra("score", 0);

        // Display completion info
        tvCompletionMessage.setText("Thank you, " + candidateName + "!");
        tvStats.setText(String.format(Locale.getDefault(),
                "Questions Answered: %d / %d", questionsAnswered, totalQuestions));

        int scorePercentage = (int) (score * 100);
        tvScoreDisplay.setText(scorePercentage + "%");

        // Set color based on score
        if (score >= 0.7) {
            tvScoreDisplay.setTextColor(getResources().getColor(android.R.color.holo_green_dark));
        } else if (score >= 0.5) {
            tvScoreDisplay.setTextColor(getResources().getColor(android.R.color.holo_orange_dark));
        } else {
            tvScoreDisplay.setTextColor(getResources().getColor(android.R.color.holo_red_dark));
        }

        // Register for robot lifecycle
        QiSDK.register(this, this);

        // Start new interview button
        btnNewInterview.setOnClickListener(v -> {
            Intent intent = new Intent(CompleteActivity.this, MainActivity.class);
            intent.setFlags(Intent.FLAG_ACTIVITY_CLEAR_TOP | Intent.FLAG_ACTIVITY_NEW_TASK);
            startActivity(intent);
            finish();
        });
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

        // Thank the candidate
        new Thread(() -> {
            try {
                Say say = SayBuilder.with(qiContext)
                        .withText("Your interview has been recorded and will be reviewed by our HR team. We wish you the best of luck!")
                        .build();
                say.run();
            } catch (Exception e) {
                Log.e(TAG, "Error speaking: " + e.getMessage());
            }
        }).start();
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
}
