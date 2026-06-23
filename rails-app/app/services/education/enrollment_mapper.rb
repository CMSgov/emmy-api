module Education
  class EnrollmentMapper
    def self.translate_nsc_response(nsc_request, nsc_resp, duration, version: :v0)
      if version != :v1 && (is_nsc_no_hit?(nsc_resp) || is_nsc_not_currently_enrolled?(nsc_resp))
        raise Education::NotFoundError, "Not Found"
      end

      status = resolve_enrollment_status(nsc_resp)
      if version != :v1 && !status
        raise "NSC response missing enrollment status"
      end

      if version == :v1
        map_v1(nsc_request, nsc_resp, status, duration)
      else
        map_v0(nsc_request, nsc_resp, status, duration)
      end
    end

    private

    def self.map_v0(nsc_request, nsc_resp, status, duration)
      details = []
      (nsc_resp["enrollmentDetails"] || []).each do |d|
        (d["enrollmentData"] || []).each do |ed|
          norm_status = EnrollmentStatus.normalize(ed["enrollmentStatus"])
          next unless norm_status

          details << {
            schoolName: d["officialSchoolName"],
            schoolCode: d["schoolCode"],
            termBeginDate: ed["termBeginDate"],
            termEndDate: ed["termEndDate"],
            enrollmentStatus: norm_status
          }
        end
      end

      Education::EnrollmentResponseV0.new(
        enrollmentStatus: status,
        enrollmentDetails: details,
        rawData: nsc_resp,
        dataSource: "NSC",
        metadata: build_metadata(duration)
      )
    end

    def self.map_v1(nsc_request, nsc_resp, status, duration)
      details = []
      (nsc_resp["enrollmentDetails"] || []).each do |d|
        next_enrollment_detail = {
          officialSchoolName: d["officialSchoolName"],
          schoolCode: d["schoolCode"],
          branchCode: d["branchCode"] || "00" # doesn't actually exist in NSC but supposedly well get once contract sign
        }
        enrollment_data = (d["enrollmentData"] || []).map do |ed|
          {
            termBeginDate: ed["termBeginDate"],
            termEndDate: ed["termEndDate"],
            enrollmentStatusCode: ed["enrollmentStatus"],
            schoolCertifiedOnDate: ed["schoolCertifiedOnDate"],
            anticipatedGraduationDate: ed["anticipatedGraduationDate"]
          }
        end
        next_enrollment_detail[:enrollmentData] = enrollment_data
        details << next_enrollment_detail
      end

      student_info = {
        personGivenName: nsc_request.is_a?(Hash) ? nsc_request.dig("studentInfoProvided", "personGivenName") : nsc_request.person_given_name,
        personSurName: nsc_request.is_a?(Hash) ? nsc_request.dig("studentInfoProvided", "personSurName") : nsc_request.person_sur_name,
        previousNames: nsc_request.is_a?(Hash) ? (nsc_request.dig("studentInfoProvided", "previousNames") || []) : nsc_request.previous_names,
        personBirthDate: nsc_request.is_a?(Hash) ? nsc_request.dig("studentInfoProvided", "personBirthDate") : nsc_request.person_birth_date
      }

      Education::EnrollmentResponseV1.new(
        enrollmentDetails: details,
        studentInfoProvided: student_info,
        transactionDetails: nsc_resp.dig("transactionDetails")
      )
    end

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
