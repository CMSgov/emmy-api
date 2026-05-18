module Education
  class EnrollmentMapper
    def self.translate_nsc_response(nsc_resp, duration)
      if is_nsc_no_hit?(nsc_resp) || is_nsc_not_currently_enrolled?(nsc_resp)
        raise Education::NotFoundError, "Not Found"
      end

      status = resolve_enrollment_status(nsc_resp)
      raise "NSC response missing enrollment status" unless status

      details = []
      (nsc_resp["enrollmentDetails"] || []).each do |d|
        (d["enrollmentData"] || []).each do |ed|
          norm_status = EnrollmentStatus.normalize(ed["enrollmentStatus"])
          next unless norm_status

          details << {
            schoolName: d["officialSchoolName"],
            termBeginDate: ed["termBeginDate"],
            termEndDate: ed["termEndDate"],
            enrollmentStatus: norm_status
          }
        end
      end

      Education::EnrollmentResponse.new(
        enrollmentStatus: status,
        enrollmentDetails: details,
        rawData: nsc_resp,
        dataSource: "NSC",
        metadata: build_metadata(duration)
      )
    end

    private

    def self.build_metadata(duration)
      {
        apiVersion: ENV["SERVICE_VERSION"] || "1.3.0",
        environment: ENV["ENVIRONMENT"] || "development",
        requestTimestamp: Time.now.utc.iso8601(3),
        responseTimestamp: Time.now.utc.iso8601(3),
        transactionId: SecureRandom.uuid,
        datasourceDurationMillis: duration.to_i
      }
    end

    def self.is_nsc_no_hit?(resp)
      hit = resp.dig("transactionDetails", "nscHit").to_s.upcase.strip
      [ "N", "NO", "FALSE", "0" ].include?(hit)
    end

    def self.is_nsc_not_currently_enrolled?(resp)
      details = resp["enrollmentDetails"] || []
      return false if details.empty?

      details.any? { |d| d["currentEnrollmentStatus"].to_s.strip.casecmp?("CN") }
    end

    def self.resolve_enrollment_status(resp)
      best = nil

      details = resp["enrollmentDetails"] || []
      details.each do |detail|
        (detail["enrollmentData"] || []).each do |item|
          status = EnrollmentStatus.normalize(item["enrollmentStatus"])
          if status && EnrollmentStatus.rank(status) > EnrollmentStatus.rank(best)
            best = status
          end
        end

        status = EnrollmentStatus.normalize_current(detail["currentEnrollmentStatus"])
        if status && EnrollmentStatus.rank(status) > EnrollmentStatus.rank(best)
          best = status
        end
      end

      return best if best

      return EnrollmentStatus::UNKNOWN_CREDIT_TIMING if is_nsc_positive_hit?(resp)

      nil
    end

    def self.is_nsc_positive_hit?(resp)
      hit = resp.dig("transactionDetails", "nscHit").to_s.upcase.strip
      return true if [ "Y", "YES", "TRUE", "1" ].include?(hit)

      resp.dig("status", "code").to_s.strip == "0"
    end
  end
end
