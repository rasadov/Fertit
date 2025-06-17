-- Create trigger for setting timestamps automatically
CREATE TRIGGER set_subscriber_created_at
BEFORE INSERT ON subscribers
FOR EACH ROW
EXECUTE FUNCTION set_created_at();

CREATE TRIGGER set_login_attempt_timestamp
BEFORE INSERT ON login_attempt
FOR EACH ROW
EXECUTE FUNCTION set_created_at();

-- Function to check and handle unsubscriptions
CREATE OR REPLACE FUNCTION check_subscription_status()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.incidents = FALSE AND NEW.new_features = FALSE AND NEW.others = FALSE AND NEW.news = FALSE AND NEW.policy_updates = FALSE THEN
        -- All preferences are turned off, we could log this or handle it differently
        -- For now, just updating a tracking field would be one option
        RAISE NOTICE 'User % has unsubscribed from all categories', NEW.email;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to monitor subscription status
CREATE TRIGGER check_subscriber_status
AFTER UPDATE ON subscribers
FOR EACH ROW
EXECUTE FUNCTION check_subscription_status();

-- Function to validate an email format before insertion
CREATE OR REPLACE FUNCTION validate_email()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.email !~ '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$' THEN
        RAISE EXCEPTION 'Invalid email format: %', NEW.email;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to validate email
CREATE TRIGGER validate_subscriber_email
BEFORE INSERT OR UPDATE ON subscribers
FOR EACH ROW
EXECUTE FUNCTION validate_email();