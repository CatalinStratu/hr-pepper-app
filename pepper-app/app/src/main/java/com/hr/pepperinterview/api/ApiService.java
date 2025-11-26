package com.hr.pepperinterview.api;

import com.hr.pepperinterview.models.Interview;
import com.hr.pepperinterview.models.Response;

import retrofit2.Call;
import retrofit2.http.Body;
import retrofit2.http.GET;
import retrofit2.http.POST;
import retrofit2.http.PUT;
import retrofit2.http.Path;

public interface ApiService {

    // Create a new interview
    @POST("/api/interviews")
    Call<Interview.CreateInterviewResponse> createInterview(@Body Interview.CreateInterviewRequest request);

    // Update interview status
    @PUT("/api/interviews/{id}")
    Call<Void> updateInterview(@Path("id") String id, @Body Interview.UpdateInterviewRequest request);

    // Submit a response
    @POST("/api/responses")
    Call<Response.SubmitResponseResult> submitResponse(@Body Response.SubmitResponseRequest request);

    // Get interview details
    @GET("/api/interviews/{id}")
    Call<Interview> getInterview(@Path("id") String id);
}
