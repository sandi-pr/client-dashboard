function shortenType(type) {
    // Remove vowels and take the first three consonants
    return type.replace(/[aeiou]/gi, '').toUpperCase().slice(0, 3);
}

function padNumber(input, length) {
    return input.toString().padStart(length, '0');
}

export function generateUniqueId(input, type) {
    const paddedNumber = padNumber(input, 3); // Ensure the number is at least 3 digits
    const shortType = shortenType(type);

    // Combine the shortened type and padded number
    const result = `${shortType}${paddedNumber}`;

    return result;
}

export function decodeUniqueId(uniqueId) {
    const shortType = uniqueId.match(/[A-Z]{1,3}/)[0];
    const paddedNumber = uniqueId.slice(shortType.length);

    const id = parseInt(paddedNumber, 10);
    return { type: shortType, id };
}
