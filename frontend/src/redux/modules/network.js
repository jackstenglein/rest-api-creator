
export const STATUS_PENDING = "FETCHING";
export const STATUS_SUCCESS = "SUCCESS";
export const STATUS_FAILURE = "FAILURE";
export const STATUS_NONE = "NONE";


export function pending() {
  return {
    status: STATUS_PENDING
  };
}

export function success() {
  return {
    status: STATUS_SUCCESS
  };
}

export function failure(error) {
  return {
    status: STATUS_FAILURE,
    error: error
  };
}

export function none() {
  return {
    status: STATUS_NONE
  };
}
