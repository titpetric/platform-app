package opa

# Features block (YAML → JSON)
features := {"allow_registration":true,"email_reverify_days":180,"mfa":true,"password_expiry_days":30}

# Flow definitions
flows := {
  "forgotten_password": {
    "steps": {
      "forgot_password": {
        "view": "forgot_password.vuego",
        "link": "/forgot-password",
        "next": {"error":"forgot_error","ok":"send_reset"}
      },
      "send_reset": {
        "view": "reset_sent.vuego",
        "link": "/forgot-password/sent",
        "next": {"ok":"reset_password"}
      },
      "reset_password": {
        "view": "reset_password.vuego",
        "link": "/reset-password",
        "next": {"error":"reset_password_error","ok":"success"}
      },
      "success": {
        "view": "",
        "link": "/login/success",
        "next": null
      }
    }
  },
  "login": {
    "steps": {
      "login": {
        "view": "login.vuego",
        "link": "/login",
        "next": {"error":"login_error","ok":"check_mfa"}
      },
      "check_mfa": {
        "view": "mfa_challenge.vuego",
        "link": "/mfa",
        "enabled_if": "features.mfa",
        "next": {"error":"check_mfa_error","ok":"success"}
      },
      "success": {
        "view": "",
        "link": "/login/success",
        "next": null
      }
    }
  },
  "policies": {
    "steps": {
      "password_expired": {
        "view": "password_expired.vuego",
        "link": "/password-expired",
        "enabled_if": "features.password_expiry_days > 0",
        "next": {"error":"reset_password_error","ok":"success"}
      },
      "email_reverify": {
        "view": "email_reverify.vuego",
        "link": "/verify-email",
        "enabled_if": "features.email_reverify_days > 0",
        "next": {"error":"email_reverify_error","ok":"success"}
      },
      "success": {
        "view": "",
        "link": "/login/success",
        "next": null
      }
    }
  },
  "registration": {
    "steps": {
      "register_form": {
        "view": "register.vuego",
        "link": "/register",
        "next": {"error":"register_error","ok":"email_verify"}
      },
      "email_verify": {
        "view": "email_verify.vuego",
        "link": "/verify-email",
        "next": {"error":"email_verify_error","ok":"success"}
      },
      "success": {
        "view": "",
        "link": "/login/success",
        "next": null
      }
    }
  }
}

# Check if a step is enabled
step_enabled(flow, step) {
    not flows[flow].steps[step].enabled_if
}

step_enabled(flow, step) {
    cond := flows[flow].steps[step].enabled_if
    data.userflow.evaluate_condition[cond]
}

# Evaluate conditions (simple features mapping)
evaluate_condition[k] = v {
    parts := split(k, ".")
    v := features[parts[1]]
}

# Determine next step
next_step[step] {
    f := input.flow
    s := input.step
    result := input.result

    step_enabled(f, s)
    step := flows[f].steps[s].next[result]
}

# Get the template/view for the next step
view_for_step[view] {
    f := input.flow
    s := next_step[_]
    view := flows[f].steps[s].view
}

# Optional: get link for the next step
link_for_step[link] {
    f := input.flow
    s := next_step[_]
    link := flows[f].steps[s].link
}
