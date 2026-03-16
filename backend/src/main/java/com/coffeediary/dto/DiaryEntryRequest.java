package com.coffeediary.dto;

import jakarta.validation.constraints.Max;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotNull;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class DiaryEntryRequest {

    @NotNull(message = "Date/time is required")
    private LocalDateTime dateTime;

    private Long sieveId;

    @Builder.Default
    private Integer temperature = 93;

    private Long coffeeId;

    private Double grindSize;

    private Double inputWeight;

    private Double outputWeight;

    private Integer timeSeconds;

    @Min(value = 1, message = "Rating must be between 1 and 5")
    @Max(value = 5, message = "Rating must be between 1 and 5")
    private Integer rating;

    private String notes;
}
