module Education
  module EnrollmentStatus
    FULL_TIME = "FULL_TIME"
    THREE_QUARTERS_TIME = "THREE_QUARTERS_TIME"
    HALF_TIME = "HALF_TIME"
    LESS_THAN_HALF_TIME = "LESS_THAN_HALF_TIME"
    UNKNOWN_CREDIT_TIMING = "ENROLLMENT_STATUS_UNKNOWN_CREDIT_TIMING"

    RANKS = {
      FULL_TIME => 5,
      THREE_QUARTERS_TIME => 4,
      HALF_TIME => 3,
      LESS_THAN_HALF_TIME => 2,
      UNKNOWN_CREDIT_TIMING => 1
    }.freeze

    def self.rank(status)
      RANKS.fetch(status, 0)
    end

    def self.normalize(value)
      return nil if value.blank?

      normalized = value.to_s.upcase.strip.gsub(/[- ]/, '_')

      case normalized
      when "FULL_TIME", "F" then FULL_TIME
      when "THREE_QUARTERS_TIME", "Q", "THREE_QUARTER_TIME" then THREE_QUARTERS_TIME
      when "HALF_TIME", "H" then HALF_TIME
      when "LESS_THAN_HALF_TIME", "L" then LESS_THAN_HALF_TIME
      when "ENROLLMENT_STATUS_UNKNOWN_CREDIT_TIMING", "Y" then UNKNOWN_CREDIT_TIMING
      else nil
      end
    end

    def self.normalize_current(value)
      normalized = value.to_s.upcase.strip
      case normalized
      when "CC" then UNKNOWN_CREDIT_TIMING
      else nil
      end
    end

    def self.all
      [FULL_TIME, THREE_QUARTERS_TIME, HALF_TIME, LESS_THAN_HALF_TIME, UNKNOWN_CREDIT_TIMING]
    end
  end
end
