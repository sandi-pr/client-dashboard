// @ts-nocheck

function writeAble(value) {

    // track of subscribers
    const subscribers = new Set();

    // add a value to the subscribers
    function subscribe(subscriber) {
        subscribers.add(subscriber);


        return function unsubscribe() {
            subscribers.delete(subscriber);
        }
    }


    // notify the subscribers

    function update(updater) {
        set(updater(value));
    }

    function set(newValue) {
        value = newValue
        subscribers.forEach((subscriber) => subscriber(value));
    }

    return {
        subscribe, set, update
    }
}

export const timelineStep = writeAble(1);


