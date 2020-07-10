CREATE OR REPLACE FUNCTION public.increase_booked() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
    BEGIN
        UPDATE slots
            SET "booked" = "booked" + 1
        WHERE slots.id = NEW.slot_id;
        RETURN NEW;
    end
    $$;


CREATE TRIGGER "increaseNumberBookedSlots" AFTER INSERT ON public.bookings FOR EACH ROW EXECUTE PROCEDURE public.increase_booked();
