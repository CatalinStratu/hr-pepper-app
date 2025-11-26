package com.hr.pepperinterview.models;

import com.google.gson.annotations.SerializedName;
import java.util.List;

public class Interview {
    @SerializedName("id")
    private String id;

    @SerializedName("candidate_name")
    private String candidateName;

    @SerializedName("candidate_email")
    private String candidateEmail;

    @SerializedName("position")
    private String position;

    @SerializedName("status")
    private String status;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getCandidateName() {
        return candidateName;
    }

    public void setCandidateName(String candidateName) {
        this.candidateName = candidateName;
    }

    public String getCandidateEmail() {
        return candidateEmail;
    }

    public void setCandidateEmail(String candidateEmail) {
        this.candidateEmail = candidateEmail;
    }

    public String getPosition() {
        return position;
    }

    public void setPosition(String position) {
        this.position = position;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    // Response wrapper classes
    public static class CreateInterviewResponse {
        @SerializedName("interview")
        private Interview interview;

        @SerializedName("questions")
        private List<Question> questions;

        public Interview getInterview() {
            return interview;
        }

        public List<Question> getQuestions() {
            return questions;
        }
    }

    public static class CreateInterviewRequest {
        @SerializedName("candidate_name")
        private String candidateName;

        @SerializedName("candidate_email")
        private String candidateEmail;

        @SerializedName("position")
        private String position;

        public CreateInterviewRequest(String candidateName, String candidateEmail, String position) {
            this.candidateName = candidateName;
            this.candidateEmail = candidateEmail;
            this.position = position;
        }
    }

    public static class UpdateInterviewRequest {
        @SerializedName("status")
        private String status;

        @SerializedName("score")
        private Double score;

        public UpdateInterviewRequest(String status, Double score) {
            this.status = status;
            this.score = score;
        }
    }
}
