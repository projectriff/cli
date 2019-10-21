package commands

const bash_completion_func = `
__riff_override_flag_list=(--kubeconfig --namespace -n)
__riff_override_flags()
{
    local ${__riff_override_flag_list[*]##*-} two_word_of of var
    for w in "${words[@]}"; do
        if [ -n "${two_word_of}" ]; then
            eval "${two_word_of##*-}=\"${two_word_of}=\${w}\""
            two_word_of=
            continue
        fi
        for of in "${__riff_override_flag_list[@]}"; do
            case "${w}" in
                ${of}=*)
                    eval "${of##*-}=\"${w}\""
                    ;;
                ${of})
                    two_word_of="${of}"
                    ;;
            esac
        done
    done
    for var in "${__riff_override_flag_list[@]##*-}"; do
        if eval "test -n \"\$${var}\""; then
            eval "echo -n \${${var}}' '"
        fi
    done
}

__riff_list_namespaces()
{
    local template
    template="{{ range .items  }}{{ .metadata.name }} {{ end }}"
    local kubectl_out
    # TODO decouple from kubectl
    if kubectl_out=$(kubectl get $(__riff_override_flags) -o template --template="${template}" namespace 2>/dev/null); then
        COMPREPLY=( $( compgen -W "${kubectl_out}[*]" -- "$cur" ) )
    fi
}

__riff_list_knative_configurations()
{
    local template
    template="{{ range .items  }}{{ .metadata.name }} {{ end }}"
    local kubectl_out
    # TODO decouple from kubectl
    if kubectl_out=$(kubectl get $(__riff_override_flags) -o template --template="${template}" configurations.serving.knative.dev 2>/dev/null); then
        COMPREPLY=( $( compgen -W "${kubectl_out}[*]" -- "$cur" ) )
    fi
}

__riff_list_knative_services()
{
    local template
    template="{{ range .items  }}{{ .metadata.name }} {{ end }}"
    local kubectl_out
    # TODO decouple from kubectl
    if kubectl_out=$(kubectl get $(__riff_override_flags) -o template --template="${template}" services.serving.knative.dev 2>/dev/null); then
        COMPREPLY=( $( compgen -W "${kubectl_out}[*]" -- "$cur" ) )
    fi
}

__riff_list_streaming_provisioner_services()
{
    local template
    template="{{ range .items  }}{{ .metadata.name }} {{ end }}"
    local kubectl_out
    # TODO decouple from kubectl
    if kubectl_out=$(kubectl get $(__riff_override_flags) -o template --template="${template}" --selector streaming.projectriff.io/provisioner services 2>/dev/null); then
        COMPREPLY=( $( compgen -W "${kubectl_out}[*]" -- "$cur" ) )
    fi
}

__riff_list_functions()
{
	__riff_list_resource 'function list'
}

__riff_list_containers()
{
	__riff_list_resource 'container list'
}

__riff_list_applications()
{
	__riff_list_resource 'application list'
}

__riff_list_resource()
{
	__riff_debug "listing $1"
    local riff_output out
    if riff_output=$(riff $1 $(__riff_override_flags) 2>/dev/null); then
        out=($(echo "${riff_output}" | awk 'NR>1 {print $1}'))
        COMPREPLY=( $( compgen -W "${out[*]}" -- "$cur" ) )
    fi
}

__riff_custom_func() {
    case ${last_command} in
        riff_application_delete | riff_application_status | riff_application_tail)
            __riff_list_resource 'application list'
            return
            ;;
        riff_container_delete | riff_container_status)
            __riff_list_resource 'container list'
            return
            ;;
        riff_core_deployer_delete | riff_core_deployer_status | riff_core_deployer_tail)
            __riff_list_resource 'core deployer list'
            return
            ;;
        riff_crediential_delete)
            __riff_list_resource 'credential list'
            return
            ;;
        riff_function_delete | riff_function_status | riff_function_tail)
            __riff_list_resource 'function list'
            return
            ;;
        riff_knative_deployer_delete | riff_knative_deployer_status | riff_knative_deployer_tail)
            __riff_list_resource 'knative deployer list'
            return
            ;;
        riff_knative_adapter_delete | riff_knative_adapter_status)
            __riff_list_resource 'knative adapter list'
            return
            ;;
        riff_streaming_kafka-provider_delete | riff_streaming_kafka-provider_status)
            __riff_list_resource 'streaming kafka-provider list'
            return
            ;;
        riff_streaming_processor_delete | riff_streaming_processor_status | riff_streaming_processor_tail)
            __riff_list_resource 'streaming processor list'
            return
            ;;
        riff_streaming_stream_delete | riff_streaming_stream_status)
            __riff_list_resource 'streaming stream list'
            return
            ;;
        *)
            ;;
    esac
}
`
